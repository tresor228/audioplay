function playAudio(audioItem) {
    // Récupérer les détails complets de l'audio
    fetch(`/api/v1/audio/${audioItem.id}`)
        .then(response => response.json())
        .then(details => {
            audioResults.innerHTML = `
                <div class="audio-player" style="margin: 20px; padding: 20px; background: white; border-radius: 12px;">
                    <img src="${details.thumbnail}" alt="${details.title}" style="width: 100%; border-radius: 8px; margin-bottom: 10px;">
                    <h3>${details.title}</h3>
                    <p>${details.artist}</p>
                    <audio controls autoplay style="width: 100%; margin-top: 10px;">
                        <source src="${details.stream_url}" type="audio/mpeg">
                        Your browser does not support the audio element.
                    </audio>
                </div>
            `;
        })
        .catch(error => {
            console.error('Error fetching audio details:', error);
            audioResults.innerHTML = '<div class="empty-state"><h3>Error loading audio</h3></div>';
        });
}