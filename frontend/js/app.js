document.addEventListener('DOMContentLoaded', function() {
    // Éléments DOM
    const elements = {
        searchInput: document.getElementById('searchInput'),
        searchType: document.getElementById('searchType'),
        searchButton: document.getElementById('searchButton'),
        audioSection: document.getElementById('audioSection'),
        videoSection: document.getElementById('videoSection'),
        audioResults: document.getElementById('audioResults'),
        videoResults: document.getElementById('videoResults'),
        videoPlayerContainer: document.getElementById('videoPlayerContainer'),
        videoPlayer: document.getElementById('videoPlayer'),
        tabButtons: document.querySelectorAll('.tab-btn')
    };

    // Gestion des onglets
    function handleTabSwitch(tab) {
        // Retirer la classe active de tous les boutons
        elements.tabButtons.forEach(btn => btn.classList.remove('active'));
        
        // Masquer toutes les sections
        elements.audioSection.classList.add('hidden');
        elements.videoSection.classList.add('hidden');
        elements.videoPlayerContainer.classList.add('hidden');

        // Afficher la section correspondante
        switch(tab) {
            case 'audio':
                elements.audioSection.classList.remove('hidden');
                break;
            case 'video':
                elements.videoSection.classList.remove('hidden');
                break;
            case 'favorites':
                // À implémenter: afficher les favoris
                break;
        }
    }

    // Initialisation des événements
    function initEventListeners() {
        // Gestion des onglets
        elements.tabButtons.forEach(button => {
            button.addEventListener('click', () => {
                const tab = button.getAttribute('data-tab');
                button.classList.add('active');
                handleTabSwitch(tab);
            });
        });

        // Recherche
        elements.searchButton.addEventListener('click', performSearch);
        elements.searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') performSearch();
        });
    }

    // Fonction de recherche principale
    async function performSearch() {
        const query = elements.searchInput.value.trim();
        if (!query) return;
        
        try {
            if (elements.searchType.value === 'audio') {
                await searchAudio(query);
                handleTabSwitch('audio');
            } else {
                await searchVideo(query);
                handleTabSwitch('video');
            }
        } catch (error) {
            console.error('Search error:', error);
            showErrorMessage(elements.audioResults, 'Erreur lors de la recherche');
            showErrorMessage(elements.videoResults, 'Erreur lors de la recherche');
        }
    }

    // Recherche audio
    async function searchAudio(query) {
        try {
            const response = await fetch(`/api/v1/search/audio?q=${encodeURIComponent(query)}`);
            if (!response.ok) throw new Error('Network response was not ok');
            const data = await response.json();
            displayAudioResults(data);
        } catch (error) {
            console.error('Audio search error:', error);
            throw error;
        }
    }

    // Recherche vidéo
    async function searchVideo(query) {
        try {
            const response = await fetch(`/api/v1/search/video?q=${encodeURIComponent(query)}`);
            if (!response.ok) throw new Error('Network response was not ok');
            const data = await response.json();
            displayVideoResults(data);
        } catch (error) {
            console.error('Video search error:', error);
            throw error;
        }
    }

    // Affichage des résultats audio
    function displayAudioResults(items) {
        elements.audioResults.innerHTML = '';
        
        if (!items || items.length === 0) {
            showEmptyState(elements.audioResults, 'Aucun résultat trouvé');
            return;
        }

        items.forEach(item => {
            const audioItem = document.createElement('div');
            audioItem.className = 'audio-item';
            audioItem.innerHTML = `
                <img src="${item.thumbnail}" alt="${item.title}" class="thumbnail">
                <div class="title">${item.title}</div>
            `;
            audioItem.addEventListener('click', () => playAudio(item));
            elements.audioResults.appendChild(audioItem);
        });
    }

    // Affichage des résultats vidéo
    function displayVideoResults(items) {
        elements.videoResults.innerHTML = '';
        
        if (!items || items.length === 0) {
            showEmptyState(elements.videoResults, 'Aucun résultat trouvé');
            return;
        }

        items.forEach(item => {
            const videoItem = document.createElement('div');
            videoItem.className = 'video-item';
            videoItem.innerHTML = `
                <img src="${item.thumbnail}" alt="${item.title}" class="thumbnail">
                <div class="title">${item.title}</div>
            `;
            videoItem.addEventListener('click', () => playVideo(item.id));
            elements.videoResults.appendChild(videoItem);
        });
    }

    // Lecture audio
    async function playAudio(audioItem) {
        try {
            const response = await fetch(`/api/v1/audio/${audioItem.id}`);
            if (!response.ok) throw new Error('Network response was not ok');
            const details = await response.json();
            
            elements.audioResults.innerHTML = `
                <div class="audio-player">
                    <img src="${details.thumbnail}" alt="${details.title}" class="player-thumbnail">
                    <h3>${details.title}</h3>
                    <p>${details.artist || ''}</p>
                    <audio controls autoplay class="audio-element">
                        <source src="${details.stream_url}" type="audio/mpeg">
                        Votre navigateur ne supporte pas l'élément audio.
                    </audio>
                </div>
            `;
        } catch (error) {
            console.error('Error playing audio:', error);
            showErrorMessage(elements.audioResults, 'Erreur lors du chargement de l\'audio');
        }
    }

    // Lecture vidéo
    function playVideo(videoId) {
        elements.videoPlayer.src = `https://www.youtube.com/embed/${videoId}?autoplay=1`;
        elements.videoPlayerContainer.classList.remove('hidden');
    }

    // Affichage d'un état vide
    function showEmptyState(container, message) {
        container.innerHTML = `<div class="empty-state"><h3>${message}</h3></div>`;
    }

    // Affichage d'un message d'erreur
    function showErrorMessage(container, message) {
        container.innerHTML = `<div class="error-state"><h3>${message}</h3></div>`;
    }

    // Initialisation
    initEventListeners();
});