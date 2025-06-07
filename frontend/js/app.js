document.addEventListener('DOMContentLoaded', function() {
            // Éléments DOM
            const searchInput = document.getElementById('searchInput');
            const searchType = document.getElementById('searchType');
            const searchButton = document.getElementById('searchButton');
            const audioSection = document.getElementById('audioSection');
            const videoSection = document.getElementById('videoSection');
            const favoritesSection = document.getElementById('favoritesSection');
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
                    hideAllSections();
                    
                    if (tab === 'audio') {
                        audioSection.classList.remove('hidden');
                        searchType.value = 'audio';
                    } else if (tab === 'video') {
                        videoSection.classList.remove('hidden');
                        searchType.value = 'video';
                    } else if (tab === 'favorites') {
                        favoritesSection.classList.remove('hidden');
                    }
                });
            });
            
            function hideAllSections() {
                audioSection.classList.add('hidden');
                videoSection.classList.add('hidden');
                favoritesSection.classList.add('hidden');
                videoPlayerContainer.classList.add('hidden');
            }
            
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
                    // Basculer vers l'onglet audio
                    document.querySelector('[data-tab="audio"]').click();
                } else {
                    searchVideo(query);
                    // Basculer vers l'onglet vidéo
                    document.querySelector('[data-tab="video"]').click();
                }
            }
            
            async function searchAudio(query) {
                try {
                    // Simulation de données pour la démo
                    const mockData = [
                        {
                            id: 1,
                            title: `Audio: ${query} - Résultat 1`,
                            thumbnail: `https://picsum.photos/300/180?random=audio1&query=${query}`
                        },
                        {
                            id: 2,
                            title: `Audio: ${query} - Résultat 2`,
                            thumbnail: `https://picsum.photos/300/180?random=audio2&query=${query}`
                        },
                        {
                            id: 3,
                            title: `Audio: ${query} - Résultat 3`,
                            thumbnail: `https://picsum.photos/300/180?random=audio3&query=${query}`
                        }
                    ];
                    
                    displayAudioResults(mockData);
                    
                } catch (error) {
                    console.error('Erreur lors de la recherche audio:', error);
                    audioResults.innerHTML = '<div class="empty-state"><h3>Erreur lors de la recherche</h3></div>';
                }
            }
            
            async function searchVideo(query) {
                try {
                    // Simulation de données pour la démo
                    const mockData = [
                        {
                            id: 'dQw4w9WgXcQ',
                            title: `Vidéo: ${query} - Résultat 1`,
                            thumbnail: `https://picsum.photos/300/180?random=video1&query=${query}`
                        },
                        {
                            id: 'dQw4w9WgXcQ',
                            title: `Vidéo: ${query} - Résultat 2`,
                            thumbnail: `https://picsum.photos/300/180?random=video2&query=${query}`
                        },
                        {
                            id: 'dQw4w9WgXcQ',
                            title: `Vidéo: ${query} - Résultat 3`,
                            thumbnail: `https://picsum.photos/300/180?random=video3&query=${query}`
                        }
                    ];
                    
                    displayVideoResults(mockData);
                    
                } catch (error) {
                    console.error('Erreur lors de la recherche vidéo:', error);
                    videoResults.innerHTML = '<div class="empty-state"><h3>Erreur lors de la recherche</h3></div>';
                }
            }
            
            function displayAudioResults(data) {
                audioResults.innerHTML = '';
                
                if (data.length === 0) {
                    audioResults.innerHTML = '<div class="empty-state"><h3>Aucun résultat trouvé</h3></div>';
                    return;
                }
                
                data.forEach((item, index) => {
                    const audioItem = document.createElement('div');
                    audioItem.className = 'media-item';
                    audioItem.style.animationDelay = `${index * 0.1}s`;
                    audioItem.innerHTML = `
                        <img src="${item.thumbnail}" alt="${item.title}" class="thumbnail">
                        <div class="media-title">${item.title}</div>
                    `;
                    audioItem.addEventListener('click', () => playAudio(item));
                    audioResults.appendChild(audioItem);
                });
            }
            
            function displayVideoResults(data) {
                videoResults.innerHTML = '';
                
                if (data.length === 0) {
                    videoResults.innerHTML = '<div class="empty-state"><h3>Aucun résultat trouvé</h3></div>';
                    return;
                }
                
                data.forEach((item, index) => {
                    const videoItem = document.createElement('div');
                    videoItem.className = 'media-item';
                    videoItem.style.animationDelay = `${index * 0.1}s`;
                    videoItem.innerHTML = `
                        <img src="${item.thumbnail}" alt="${item.title}" class="thumbnail">
                        <div class="media-title">${item.title}</div>
                    `;
                    videoItem.addEventListener('click', () => playVideo(item.id));
                    videoResults.appendChild(videoItem);
                });
            }
            
            function playAudio(audioItem) {
                console.log('Lecture audio:', audioItem);
                // Ici vous pouvez implémenter la lecture audio
            }
            
            function playVideo(videoId) {
                videoPlayer.src = `https://www.youtube.com/embed/${videoId}?autoplay=1`;
                videoPlayerContainer.classList.remove('hidden');
                videoPlayerContainer.scrollIntoView({ behavior: 'smooth' });
            }
            
            // Initialisation : afficher la section audio par défaut
            audioSection.classList.remove('hidden');
        });