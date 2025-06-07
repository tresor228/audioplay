document.addEventListener('DOMContentLoaded', function() {
    // Éléments DOM
    const searchInput = document.getElementById('searchInput');
    const searchType = document.getElementById('searchType');
    const searchButton = document.getElementById('searchButton');
    const audioSection = document.getElementById('audioSection');
    const videoSection = document.getElementById('videoSection');
    const audioResults = document.getElementById('audioResults');
    const videoResults = document.getElementById('videoResults');
    const videoPlayerContainer = document.getElementById('videoPlayerContainer');
    const videoPlayer = document.getElementById('videoPlayer');
    const tabButtons = document.querySelectorAll('.tab-btn');
    
    // Gestion des onglets
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            // Retirer la classe active de tous les boutons
            tabButtons.forEach(btn => btn.classList.remove('active'));
            // Ajouter la classe active au bouton cliqué
            button.classList.add('active');
            
            // Gérer l'affichage des sections
            const tab = button.getAttribute('data-tab');
            if (tab === 'audio') {
                audioSection.classList.remove('hidden');
                videoSection.classList.add('hidden');
                videoPlayerContainer.classList.add('hidden');
            } else if (tab === 'video') {
                audioSection.classList.add('hidden');
                videoSection.classList.remove('hidden');
            } else if (tab === 'favorites') {
                // À implémenter: afficher les favoris
                audioSection.classList.add('hidden');
                videoSection.classList.add('hidden');
                videoPlayerContainer.classList.add('hidden');
            }
        });
    });
    
    // Recherche
    searchButton.addEventListener('click', performSearch);
    searchInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            performSearch();
        }
    });
    
    function performSearch() {
        const query = searchInput.value.trim();
        if (query === '') return;
        
        const type = searchType.value;
        
        if (type === 'audio') {
            searchAudio(query);
        } else {
            searchVideo(query);
        }
    }
    
    async function searchAudio(query) {
        try {
            const response = await fetch(`/api/v1/search/audio?q=${encodeURIComponent(query)}`);
            const data = await response.json();
            
            audioResults.innerHTML = '';
            
            if (data.length === 0) {
                audioResults.innerHTML = '<p>Aucun résultat trouvé</p>';
                return;
            }
            
            data.forEach(item => {
                const audioItem = document.createElement('div');
                audioItem.className = 'audio-item';
                audioItem.innerHTML = `
                    <img src="${item.thumbnail}" alt="${item.title}" class="thumbnail">
                    <div class="title">${item.title}</div>
                `;
                audioItem.addEventListener('click', () => playAudio(item));
                audioResults.appendChild(audioItem);
            });
            
            // Afficher la section audio
            audioSection.classList.remove('hidden');
            videoSection.classList.add('hidden');
            videoPlayerContainer.classList.add('hidden');
            
        } catch (error) {
            console.error('Erreur lors de la recherche audio:', error);
            audioResults.innerHTML = '<p>Erreur lors de la recherche</p>';
        }
    }
    
    async function searchVideo(query) {
        try {
            const response = await fetch(`/api/v1/search/video?q=${encodeURIComponent(query)}`);
            const data = await response.json();
            
            videoResults.innerHTML = '';
            
            if (data.length === 0) {
                videoResults.innerHTML = '<p>Aucun résultat trouvé</p>';
                return;
            }
            
            data.forEach(item => {
                const videoItem = document.createElement('div');
                videoItem.className = 'video-item';
                videoItem.innerHTML = `
                    <img src="${item.thumbnail}" alt="${item.title}" class="thumbnail">
                    <div class="title">${item.title}</div>
                `;
                videoItem.addEventListener('click', () => playVideo(item.id));
                videoResults.appendChild(videoItem);
            });
            
            // Afficher la section vidéo
            audioSection.classList.add('hidden');
            videoSection.classList.remove('hidden');
            
        } catch (error) {
            console.error('Erreur lors de la recherche vidéo:', error);
            videoResults.innerHTML = '<p>Erreur lors de la recherche</p>';
        }
    }
    
    function playAudio(audioItem) {
        // À implémenter: jouer l'audio
        console.log('Lecture audio:', audioItem);
    }
    
    function playVideo(videoId) {
        videoPlayer.src = `https://www.youtube.com/embed/${videoId}?autoplay=1`;
        videoPlayerContainer.classList.remove('hidden');
    }
});