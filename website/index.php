<?
require 'Slim/Slim.php';
require 'config.php';
\Slim\Slim::registerAutoloader();

date_default_timezone_set("Europe/Paris");

$app = new \Slim\Slim();
$bdd = new PDO('mysql:host='.$bdd['host'].';dbname='.$bdd['dbname'], $bdd['user'], $bdd['passwd'], array(
        PDO::ATTR_ERRMODE => PDO::ERRMODE_WARNING));


error_reporting(E_ALL);
ini_set('display_errors', 1);

// openid
//include('/usr/share/php/openid/consumer/consumer.php');
//$consumer   =& AriseOpenID::getInstance();
//$openid_url =  !empty($_POST['openid_url']) ? $_POST['openid_url'] : NULL;
//$required = array('http://openid.iiens.net/types/identifiant');
//$consumer->authenticate($openid_url, $required);
    

// ariseid
//$consumer_key = '';
//$consumer_secret = '';
//$consumer_private_key = '';

ini_set('session.gc_maxlifetime', 3600*24*7);
session_set_cookie_params(3600*24*7);
//require_once("/usr/share/php/ariseid/client/OAuthAriseClient.php");

//$consumer = OAuthAriseClient::getInstance(
//    $consumer_key,
//    $consumer_secret,
 //   $consumer_private_key
//);

//if (isset($_POST['openid_url'])) {
//    $consumer->authenticate();
//}

//if ($consumer->has_just_authenticated()) {
//    session_regenerate_id();
//    $consumer->session_id_changed();
//}




// routes

$app->get('/fav', 'fav');
$app->get('/later/', 'later');
$app->get('/later/:nick', 'later');
$app->get('/search/:view', 'byView');
$app->get('/api/:chan', 'apiChan');
$app->get('/:chan/senders/:sender', 'bySender');
$app->get('/:chan/senders/', 'senders');
$app->get('/:chan/fav', 'fav');
$app->get('/:chan/tags/:tag', 'byTag')->name('tag');
$app->get('/:chan/tags/', 'tags');
$app->get('/:chan/widget', 'widget');
$app->get('/:chan/:date', 'day')->name('day');
$app->get('/:chan/', 'days');
$app->get('/', 'index');

$app->post('/fav', 'favPost');


function days ($chanUrl) {
    $app = \Slim\Slim::getInstance();
    $chan = '#'.$chanUrl;

    global $bdd;

    include('includes/header.php');
    include('includes/menu.php');
    echo <<<INDEXHEAD
<div class="header">Log d'activit&eacute; PlayBot</div>
<div class="content">
INDEXHEAD;



    /***************************
    * Génération du calendrier *
    ***************************/

    // décalage (mois précédent : $dif == -1)
    $dif = (isset($_GET['dif'])) ? $_GET['dif'] : 0;

    // on récupère la date actuelle;
    $year = date('Y');
    $month = date('n') + $dif;

    while ($month > 12) {
        $year++;
        $month -= 12;
    }

    while ($month < 1) {
        $year--;
        $month += 12;
    }

    $day = date('j');
    $dayWeek = date('N', mktime(0, 0, 0, $month, 1, $year)); // jour de la semaine du premier du mois

    // on récupère les jours du mois pour lesquels des liens ont été postés
    $reponse = $bdd->prepare('SELECT DISTINCT DAY(date)
        FROM playbot_chan
        WHERE MONTH(date) = '.$month.'
        AND YEAR(date) = '.$year.'
        AND chan = :chan
        ORDER BY date');
    $reponse->bindParam(':chan', $chan, PDO::PARAM_STR);
    $reponse->execute();


    // en tête du tableau (mois, année)
    echo "<table class='calendar'><thead><tr><td style='text-align: center' colspan='7'><a href='?dif=". ($dif - 1) ."'><<</a>  $month/$year  <a href='?dif=". ($dif + 1) ."'>>></a></td></tr></thead>\n"; 

    // avant de parcourir les résultats, on se postionne au bon jour de la semaine
    echo '<tr>';
    for ($i=1; $i < $dayWeek; $i++)
        echo '<td></td>';
    
    // on parcours les résultats, on enregistre le jour courant du mois pour combler les trous
    $curDay = 1;
    while ($donnees = $reponse->fetch()) {
        while ($curDay <= $donnees[0]) {
            if ($dayWeek == 1)
                echo "\n</tr>\n</tr>\n";

            if ($curDay == $donnees[0])
                echo "<td><a href='".$app->urlFor('day', array(
                        'chan' => $chanUrl,
                        'date' => "$year-$month-$donnees[0]"
                    ))."'>$donnees[0]</a></td>\n";
            else
                echo "<td>$curDay</td>";

            $curDay++;
            $dayWeek++;
            if ($dayWeek > 7)
                   $dayWeek = 1;
        }
    }


    // fin du tableau
    while ($curDay <= date('t', mktime(0, 0, 0, $month, 1, $year))) {
        if ($dayWeek == 1)
            echo "\n</tr>\n</tr>\n";

        echo "<td>$curDay</td>";

        $curDay++;
        $dayWeek++;
        if ($dayWeek > 7)
               $dayWeek = 1;
    }
    
    if ($dayWeek != 1) 
        for ($i=$dayWeek; $i <= 7; $i++)
            echo '<td></td>';

    echo "</tr></table>\n";

    /*********************
     * fin du calendrier *
     *********************/


    $nbr_senders = $bdd->prepare('SELECT sender_irc, COUNT(*) AS nb
        FROM playbot_chan
        WHERE chan = :chan
        AND sender_irc != "PlayBot"
        GROUP BY sender_irc
        ORDER BY nb
        DESC LIMIT 5');
    $nbr_senders->bindParam(':chan', $chan, PDO::PARAM_STR);
    $nbr_senders->execute();

    $nbr_types =  $bdd->prepare('SELECT type, COUNT(*) AS nb
        FROM playbot p
        JOIN playbot_chan pc ON p.id = pc.content
        WHERE chan = :chan
        GROUP BY type
        ORDER BY nb DESC');
    $nbr_types->bindParam(':chan', $chan, PDO::PARAM_STR);
    $nbr_types->execute();

    echo "<h2>Top 5 des posteurs de liens</h2>\n<ul>\n";

    while ($donnees = $nbr_senders->fetch()) {
        echo "<li><strong>$donnees[0] :</strong> $donnees[1]</li>\n";
    }

    echo "</ul>\n<h2>Top des sites</h2>\n<ul>\n";

    while ($donnees = $nbr_types->fetch()) {
        echo "<li><strong>$donnees[0] :</strong> $donnees[1]</li>\n";
    }

    echo <<<INDEXBOT
</ul>
</div>
INDEXBOT;

    $reponse->closeCursor();

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function fav ($chanUrl = '') {
    global $consumer;
    global $bdd;

    include('includes/header.php');
    include('includes/menu.php');

    if ($consumer->is_authenticated()) {
        $login = $consumer->api()->get_identifiant();
        //$login = $consumer->getSingle('http://openid.iiens.net/types/identifiant');

        // affichage des liens
        if ($chanUrl) {
            $req = $bdd->prepare('SELECT DISTINCT "osef", type, url, "osef", sender, title, p.id, broken, TRUE as fav
                FROM playbot_fav f
                JOIN playbot p ON f.id = p.id
                JOIN playbot_chan pc ON p.id = pc.content
                WHERE user = :login AND chan = :chan');
            $chan = '#'.$chanUrl;
            $req->bindParam(':chan', $chan, PDO::PARAM_STR);
        }
        else {
            $req = $bdd->prepare('SELECT "osef", type, url, "osef", sender, title, p.id, broken, TRUE as fav
                            FROM playbot_fav f
                            JOIN playbot p ON f.id = p.id
                            WHERE user = :login');
        }

        $req->bindParam(':login', $login);
        $req->execute();

        echo '<div class="header">Favoris</div>';
        $t = time();
        printLinks ($req, $chanUrl, true);
        print(time() - $t);
    

        // code pour irc
        // on regarde si un code existe déjà, sinon on en génère un
        $req = $bdd->prepare('SELECT code FROM playbot_codes WHERE user = :login');
        $req->bindParam(':login', $login);
        $req->execute();

        if ($req->rowCount())
            $code = current($req->fetch());
        else {
            $code = uniqid('PB', true);

            $req = $bdd->prepare('INSERT INTO playbot_codes (user, code) VALUES (:login, :code)');
            $req->bindParam(':login', $login);
            $req->bindParam(':code', $code);
            $req->execute();
        }

echo <<<EOF
<div class='content'>
<br/>
Pour utiliser les favoris avec PlayBot, utiliser la commande suivante : « <em>/query PlayBot $code</em> ».
</div>
EOF;
    }
    else {
        echo <<<FORM
<div class='content'>
    <form method='post' action='/links/fav'>
        Just click on "submit" : <input type='text' name='openid_url'>
        <input type='submit'>
    </form>
</div>
FORM;
    }

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function favPost () {
    global $consumer;
    global $bdd;
    $app = \Slim\Slim::getInstance();

    if (!$consumer->is_authenticated()) {
        $app->halt(500, 'User not connected');
        return;
    }

    //$login = $consumer->getSingle('http://openid.iiens.net/types/identifiant');
    $login = $consumer->api()->get_identifiant();

    // on regarde si la vidéo est déjà dans les favoris
    $req = $bdd->prepare('SELECT COUNT(*) FROM playbot_fav WHERE user = :user AND id = :id');
    $req->bindParam(':user', $login);
    $req->bindParam(':id', $_POST['id']);
    $req->execute();
    $isFav = $req->fetch();

    // si oui, on la supprime
    if ($isFav[0]) {
        $req = $bdd->prepare('DELETE FROM playbot_fav WHERE user = :user AND id = :id');
        $req->bindParam(':user', $login);
        $req->bindParam(':id', $_POST['id']);
        $req->execute();

        echo '0';
    }
    // sinon on l'ajoute
    else {
        $req = $bdd->prepare('INSERT INTO playbot_fav(user, id) VALUES(:user, :id)');
        $req->bindParam(':user', $login);
        $req->bindParam(':id', $_POST['id']);
        $req->execute();

        echo '1';
    }
}


function later ($nick = '') {
    global $consumer;
    global $bdd;

    $chanUrl = '';

    include('includes/header.php');
    include('includes/menu.php');

    if ($nick) {
        // affichage des liens
        $req = $bdd->prepare('
            SELECT DISTINCT
                \'nop\',
                type, url,
                sender_irc,
                sender, title,
                p.id,
                p.broken
            FROM
                playbot_later l
                JOIN playbot p ON l.content = p.id
            WHERE l.nick = :nick
        ');
        $req->bindParam(':nick', $nick);
        $req->execute();

        echo '<div class="header">Favoris</div>';
        printLinks ($req, $chanUrl);
    }
    else {
        echo '<p>Rendez-vous sur « later/<em>votre_pseudo_irc</em> » pour voir tous vos !later.</p>';
    }

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function widget($chanUrl) {
    global $consumer;

    /*
    if (!$consumer->is_authenticated()) {
        $consumer->authenticate();
    }
     */

    if ($consumer->has_just_authenticated()) {
        session_regenerate_id();
        $consumer->session_id_changed();
    }

    day($chanUrl, date('Y-m-d', time()), True, False);
}


function day ($chanUrl, $date, $reversed = False, $header = True) {
    global $bdd;
    $chan = '#'.$chanUrl;
    $dateMin = "$date 00:00:00";
    $dateMax = "$date 23:59:59";
    $query = 'SELECT pc.date, type, url, pc.sender_irc, sender, title, p.id, broken, GROUP_CONCAT(tag), pf.user as fav
        FROM playbot p
        LEFT OUTER JOIN playbot_tags USING (id)
        JOIN playbot_chan pc ON p.id = pc.content
        LEFT OUTER JOIN playbot_fav pf on p.id = pf.id
        WHERE pc.date
            BETWEEN :date_min
            AND :date_max
        AND chan = :chan
        GROUP BY id
        ORDER BY pc.date';

    if ($reversed) {
        $query .= ' DESC';
    }

    $req = $bdd->prepare($query);
    $req->bindParam(':date_min', $dateMin, PDO::PARAM_STR);
    $req->bindParam(':date_max', $dateMax, PDO::PARAM_STR);
    $req->bindParam(':chan', $chan, PDO::PARAM_STR);
    $req->execute();

    include('includes/header.php');

    if ($header) {
        include('includes/menu.php');
        echo <<<EOF
    <div class="header">Log d'activit&eacute; PlayBot</div>
EOF;
    }
    printLinks ($req, $chanUrl);

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function byView ($view) {
    global $bdd;
    $chanUrl = '';

    $req = $bdd->prepare('
        SELECT count(*)
        FROM information_schema.tables
        WHERE
            TABLE_TYPE LIKE "VIEW"
            AND TABLE_NAME = :view');
    $req->bindParam(':view', $view);
    $req->execute();

    include('includes/header.php');
    include('includes/menu.php');
    echo <<<EOF
    <div class="header">Log d'activit&eacute; PlayBot</div>
EOF;

    if ($req->fetch()[0] != 1) { 
        echo "<h3>Recherche inconnue ou expirée</h3>";
    }
    else {
        $req = $bdd->prepare('
            SELECT "osef", type, url, "idem", sender, title, id, 0
            FROM '.$view);
        $req->execute();

        printLinks ($req, $chanUrl);
    }

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function printLinks ($req, $chan, $selectTags = true) {
    global $consumer;

    echo '<p><a onclick="play ()" href="#">Play all</a> - <a onclick="next ()" href="#">Next</a></p>';
    echo '<div id="player"></div>';
    echo '<div class="content">';
    echo "<table>\n";
    echo "<tr class='table_header'>\n";
    echo "<td>id</td><td>Lien</td><td>Posteur</td><td>Auteur de la musique</td><td>Titre de la musique</td><td>Favoris</td><td>tags</td>\n";
    while ($donnees = $req->fetch()) {
        echo "<tr";
        if ($donnees[7]) echo " style='text-decoration:line-through;'";
        echo ">\n";
        echo "<td>".$donnees[6]."</td>\n";
        echo "<td>";
        switch ($donnees[1]) {
            case 'youtube':
                echo "<a class='content' href='$donnees[2]'><img alt='youtube' src='/links/img/yt.png' /></a>";
                break;
            case 'soundcloud':
                echo "<a href='$donnees[2]'><img alt='soundcloud' src='/links/img/sc.png' /></a>";
                break;
            case 'mixcloud':
                echo "<a href='$donnees[2]'><img alt='mixcloud' src='/links/img/mc.png' width='40px' /></a>";
                break;
            default:
                echo "<a href='$donnees[2]'>$donnees[1]</a>";
                break;
        }
        echo <<<EOF
</td>
<td>$donnees[3]</td>
<td>$donnees[4]</td>
<td>$donnees[5]</td>
EOF;

        if (array_key_exists('fav', $donnees)) {
            if ($donnees['fav'])
                echo "<td style='text-align:center'><img onClick='fav(".$donnees[6].")' id='".$donnees[6]."' src='/links/img/star-full.png' /></td>";
            else
                echo "<td style='text-align:center'><img onClick='fav(".$donnees[6].")' id='".$donnees[6]."' src='/links/img/star.png' /></td>";
        }
        else
            echo "<td style='text-align:center'><img onClick='fav(".$donnees[6].")' id='".$donnees[6]."' src='/links/img/star.png' /></td>";

        // on affiche les tags
        $tags = array();

        if ($selectTags) {
            global $bdd;
            $req2 = $bdd->prepare('SELECT tag
                FROM playbot_tags
                WHERE id = :id
                ORDER BY tag');
            $req2->bindParam(':id', $donnees[6]);
            $req2->execute();

            while ($result = $req2->fetch()) {
                $tags[] = $result[0];
            }
        }

        $first = true;
        echo '<td>';

        foreach ($tags as $tag) {
            if ($first)
                $first = false;
            else
                echo ' ';

            echo "<a href='/links/$chan/tags/$tag'>$tag</a>";
        }

        echo '</td>';
    }

    echo <<<EOF
</tr>
</table>
EOF;

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function senders ($chanUrl) {
    global $bdd;
    $chan = '#'.$chanUrl;
    $req = $bdd->prepare('SELECT DISTINCT(pc.sender_irc)
        FROM playbot p
        JOIN playbot_chan pc ON p.id = pc.content
        WHERE chan = :chan
        ORDER BY pc.sender_irc');
    $req->bindParam(':chan', $chan, PDO::PARAM_STR);
    $req->execute();

    include('includes/header.php');
    include('includes/menu.php');
    echo <<<EOF
<div class='content'>
<div class='header'>Liste des posteurs</div>
<ul>
EOF;
    echo '<p>Le regroupement des pseudos sera implémenté plus tard (kikoo Jonas !).</p>';
    echo '<ul>';
    while ($donnees = $req->fetch()) {
        echo '<li><a href="'.$donnees[0].'">'.$donnees[0]."</a></li>\n";
    }

    echo <<<FOOTER
</ul>
</div>
</body>
</html>
FOOTER;
}


function tags ($chanUrl) {
    global $bdd;
    $app = \Slim\Slim::getInstance();
    $chan = '#'.$chanUrl;

    $req = $bdd->prepare('SELECT tag, count(*) AS number
        FROM playbot_tags pt
        JOIN playbot_chan pc ON pt.id = pc.content
        WHERE chan = :chan
        GROUP BY tag
        ORDER BY tag');
    $req->bindParam(':chan', $chan, PDO::PARAM_STR);
    $req->execute();

    $min = PHP_INT_MAX;
    $max = - PHP_INT_MAX;

    while ($tag = $req->fetch()) {
        if ($tag['number'] < $min) $min = $tag['number'];
        if ($tag['number'] > $max) $max = $tag['number'];
        $tags[] = $tag;
    }
    
    include('includes/header.php');
    include('includes/menu.php');

    echo <<<EOF
<div class='content'>
<div class='header'>Liste des tags</div>
<div class='tags'>
EOF;

    if (!$tags) {
        echo '<p>Y a pas grand chose :(</p>';
        return;
    }

    $min_size = 10;
    $max_size = 100;

    foreach ($tags as $tag) {
        if ($max - $min != 0)
            $tag['size'] = intval($min_size + (($tag['number'] - $min) * (($max_size - $min_size) / ($max - $min))));
        else
            $tag['size'] = $max_size / 2;
        $tags_extended[] = $tag;
    }


    foreach ($tags_extended as $tag) {
        echo '<a style="font-size: '.$tag['size'].'px" href="';
        echo $app->urlFor('tag', array(
            'chan' => $chanUrl,
            'tag' => $tag[0]));
        echo '"">'.$tag[0].'</a> ';
    }

    echo <<<FOOTER
</div>
</div>
</body>
</html>
FOOTER;
}


function bySender ($chanUrl, $sender) {
    global $bdd;
    $chan = '#'.$chanUrl;
    $req = $bdd->prepare('
        SELECT pc.date, type, url, pc.sender_irc, sender, title, p.id, broken, GROUP_CONCAT(tag)
        FROM playbot p
        LEFT OUTER JOIN playbot_tags USING(id)
        JOIN playbot_chan pc ON p.id = pc.content
        WHERE pc.sender_irc = :sender
        AND chan = :chan
        GROUP BY id');
    $req->bindParam(':sender', $sender, PDO::PARAM_STR);
    $req->bindParam(':chan', $chan, PDO::PARAM_STR);
    $req->execute();

    include('includes/header.php');
    include('includes/menu.php');
    printLinks ($req, $chanUrl);

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function byTag ($chanUrl, $tag) {
    global $bdd;
    $chan = '#'.$chanUrl;

    $req = $bdd->prepare('
        SELECT pc.date, type, url, pc.sender_irc, sender, title, p.id, broken, pf.user as fav
        FROM playbot p
        NATURAL JOIN playbot_tags
        JOIN playbot_chan pc ON p.id = pc.content
        LEFT OUTER JOIN playbot_fav pf on p.id = pf.id
        WHERE chan = :chan
        AND tag = :tag
        GROUP BY id');

    $req->bindParam(':tag', $tag, PDO::PARAM_STR);
    $req->bindParam(':chan', $chan, PDO::PARAM_STR);

    $req->execute();


    include('includes/header.php');
    include('includes/menu.php');
    printLinks ($req, $chanUrl, true);

    echo <<<FOOTER
</body>
</html>
FOOTER;
}


function index () {
    global $bdd;

    $req = $bdd->prepare('
        SELECT chan
        FROM playbot_chan
        WHERE chan LIKE "#%"
        GROUP BY chan');
    $req->execute();
    
    include('includes/header.php');

    echo '<ul>';
    while ($chan = $req->fetch())
        echo "<li><a href='".substr($chan[0],1)."'>$chan[0]</a></li>";

    echo <<<FOOTER
</ul>
</body>
</html>
FOOTER;
}


function apiChan ($chanUrl) {
    global $bdd;
    $chan = '#'.$chanUrl;

    $req = $bdd->prepare('
        SELECT pc.date as date, type, url, pc.sender_irc as sender_irc, sender, title, p.id as id, broken, GROUP_CONCAT(DISTINCT tag) as tags
        FROM playbot p
        LEFT OUTER JOIN playbot_tags USING (id)
        JOIN playbot_chan pc ON p.id = pc.content
        WHERE chan = :chan
        GROUP BY id
        ORDER BY pc.date DESC
        LIMIT 10
    ');
    $req->bindParam(':chan', $chan, PDO::PARAM_STR);
    $req->execute();

    $result = $req->fetchAll(PDO::FETCH_ASSOC);
    echo json_encode($result);
}


$app->run();

?>
