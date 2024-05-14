function fav(id) {
	jQuery.ajax({
		type: 'POST', // Le type de ma requete
		url: '/fav', // L'url vers laquelle la requete sera envoyee
		data: {
			id: id // Les donnees que l'on souhaite envoyer au serveur au format JSON
		}, 
		success: function(data, textStatus, jqXHR) {
			if (data == 1)
				$('#' + id).attr('src', '/img/star-full.png');
			else
				$('#' + id).attr('src', '/img/star.png');
		},
		error: function(jqXHR, textStatus, errorThrown) {
			alert('Erreur. Vérifiez que vous êtes bien connecté.');
		}
	});
}

var indexPlay = 0;
var player;

function play() {
	url = $($('a.content')[indexPlay]).attr('href');
	if (url.indexOf("youtube") > -1) {
		id = url.substring(url.lastIndexOf("=") + 1);
		player = new YT.Player('player', {
			height: '390',
			width: '640',
			videoId: id,
			events: {
				'onReady': onPlayerReady,
				'onStateChange': onPlayerStateChange
			}
		});
	}
}

function next() {
	if (player) {
		player.destroy();
	}
	indexPlay++;
	play();
}

// autoplay video
function onPlayerReady(event) {
	event.target.playVideo();
}

// when video ends
function onPlayerStateChange(event) {
	if(event.data === 0) { 
		next();
	}
}
