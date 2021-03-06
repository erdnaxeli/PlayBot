# PlayBot


## Description

PlayBot est un bot IRC permettant de manipuler des liens audio :
+ enregistrement des liens dans une base de données ;
+ récupération et publication sur le chan des informations (titre, posteur) ;
+ normalisation du lien ;
+ sauvegarde de liens en favoris ;
+ rappel du lien en query ;
+ recherche dans la base de données.


## Installation


### Dépendances

PlayBot est codé en Perl. Il dépend des modules suivant :
+ DBI ;
+ DBD::mysql ;
+ JSON ;
+ LWP::UserAgent ;
+ HTML::Parser ;
+ HTML::Entities ;
+ POE ;
+ POE::Component::IRC ;
+ Tie::File ;
+ Moose.
+ Inline::Python

Il utilise de plus la librairie Python [Weboob](http://weboob.org).


### Configuration

Pour l'instant, seule la configuration de la base de donnée est externalisée. Vous devez créer un fichier *playbot.json* dans le même dossier que le fichier *PlayBot.pl* contenant les attributs suivant (les noms devraient être assez explicites) :
+ bdd ;
+ host ;
+ user ;
+ passwd.

Le reste de la configuration s'effectue encore directement dans le fichier *PlayBot.pl*. Les variables suivantes sont modifiables :
+ $serveur ;
+ $port ;
+ $nick ;
+ $ircname ;
+ $username ;
+ @channels ;
+ $admin ;
+ $baseurl.

Les noms des variables sont plutôt explicites. *$admin* contient le nick du l'utilisateur qui aura le droit d'administrer le bot (attention, aucune vérification autre que le nick n'est pour l'instant faîte, il est recommandé de l'enregistrer). *$baseurl* contient l'url racine du site permettant de consulter le contenu de la base de données.

Les schémas des tables nécéssaire pour la base de données sont dans le fichier *bdd.sql*.


## Utilisation

Dès qu'un lien est posté, le bot récupère le nom du contenu et de la personne l'ayant posté, et l'insère dans la base de données (grâce à l'unicité du lien et à sa normalisation, un même contenu ne peut être enregistré deux fois). Il poste ensuite sur le channel les informations récupérées avec l'id du contenu.

Lors du poste du lien, on peut ajouter des tags en les préfixant par le symbole « # ».

Les différentes commandes utilisables sont les suivantes sur le channel où se trouve le PlayBot :
+ !fav [id] : enregistre en favoris le contenu possédant l'identifiant *id*. Si l'identifiant n'est pas préciser, le dernier contenu posté (et non inséré) est utilisé.
+ !later [id [in [Xs|Xm|Xs]]] : demande au PlayBot de rappeler en query un contenu. La durée par défaut est de 6h. Si l'identifant n'est pas précisé, le dernier contenu posté est utilisé.
+ !tag [id] #tag1 #tag2 … : permet de taguer un contenu.
+ !get [-a|-all] <query> : permet de chercher un contenu. La query peut être composée de tags ou de mots (alors cherché dans le titre et l'auteur du contenu). Sans l'option « -a », !get ne recherche que parmis les contenus postés sur le channel courant. En query l'option « -a » est implicite.
