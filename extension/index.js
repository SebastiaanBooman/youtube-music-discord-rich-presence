// Extension observes the playstate element in the DOM
// If it changed, video metadata is sent to the back-end to create a rich presence on Discord

// Without mediaSession in the DOM the script cannot fetch metadata
// Alternative is grabbing metadata from HTML elements, but the script is scuffed enough as is
if (!('mediaSession' in navigator)) {
	throw new Error("No mediaSession present in navigator")
}
const mediaSession = navigator.mediaSession;

const mutationObserver = new MutationObserver(entries => {
	var timeData = document.querySelector(".time-info.style-scope.ytmusic-player-bar")
		.textContent.trim();

	let playing = mediaSession.playbackState === "playing" ? true : false

	postSongData(
		playing,
		mediaSession.metadata.artist,
		mediaSession.metadata.title,
		mediaSession.metadata.album,
		mediaSession.metadata.artwork[0].src,
		timeData
	);
});

async function postSongData(playbackState, artist, title, album, imageUrl, timeData){
	let request = new Request("http://127.0.0.1:8080/song-data", {
		method: "POST",
		body: JSON.stringify(
			{ 
				playing: playbackState,
				artist: artist,
				title: title,
				album: album,
				imageUrl: imageUrl,
				timeData: timeData
			}),
		headers: { "Content-Type": "application/json" }
	});
	fetch(request).then(response1 => {
	}).catch(error => {
		console.error(error);
	});
}

const playStateElement = document.querySelector(".play-pause-button.style-scope.ytmusic-player-bar");

mutationObserver.observe(playStateElement, {
	attributes : true,
	childList : false,
	subtree : false,
	attributeFilter: ["title"]
});
