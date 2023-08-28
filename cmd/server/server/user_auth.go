package server

import (
	"fmt"
	"log"

	pb "github.com/erdnaxeli/PlayBot/cmd/server/rpc"
)

func (s *server) handleUserAuth(msg *pb.TextMessage) (*pb.Result, error) {
	// user authentication
	s.ctcMutex.Lock()
	defer s.ctcMutex.Unlock()

	s.codesToCheck[msg.PersonName] = msg.Msg

	return makeResult(
		&pb.IrcMessage{
			Msg: "Vérification en cours…",
			To:  msg.PersonName,
		},
		&pb.IrcMessage{
			Msg: fmt.Sprintf("info %s", msg.PersonName),
			To:  "NickServ",
		},
	), nil
}

func (s *server) handleUserAuthCallback(msg *pb.TextMessage) (*pb.Result, error) {
	if len(s.codesToCheck) == 0 {
		return emptyResult, nil
	}

	log.Print("Received a message from NickServ.")
	s.ctcMutex.Lock()
	defer s.ctcMutex.Unlock()

	for nick, code := range s.codesToCheck {
		log.Printf("Trying to auth %s.", nick)

		if msg.Msg == fmt.Sprintf("Le pseudo %s%s%s n'est pas enregistré.", string(STX), nick, string(STX)) {
			log.Print("Unregistered nick")
			return makeResult(&pb.IrcMessage{
				Msg: "Il faut que ton pseudo soit enregistré auprès de NickServ pour pouvoir t'authentifier.",
				To:  nick,
			}), nil
		} else if msg.Msg != fmt.Sprintf("%s est actuellement connecté.", nick) {
			continue
		}

		log.Printf("Ok, authenticating nick %s.", nick)
		delete(s.codesToCheck, nick)

		user, err := s.repository.GetUserFromCode(code)
		if err != nil {
			return emptyResult, err
		} else if user == "" {
			log.Printf("Unknown code.")
			return makeResult(&pb.IrcMessage{
				Msg: "Code inconnu. Va sur http://nightiies.iiens.net/links/fav pour obtenir ton code personel.",
				To:  nick,
			}), nil
		}

		log.Printf("Code ok.")
		err = s.repository.SaveAssociation(user, nick)
		if err != nil {
			return emptyResult, err
		}

		return makeResult(&pb.IrcMessage{
			Msg: "Association effectuée. Utilise la commande !fav pour enregistrer un lien dans tes favoris.",
			To:  nick,
		}), nil
	}

	return emptyResult, nil
}
