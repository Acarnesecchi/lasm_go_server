// Import and configure the Firebase SDK
// These scripts are made available when the app is served or deployed on Firebase Hosting
// If not using Firebase Hosting, include the Firebase SDK for Cloud Messaging
importScripts('https://www.gstatic.com/firebasejs/8.10.1/firebase-app.js');
importScripts('https://www.gstatic.com/firebasejs/8.10.1/firebase-messaging.js');

firebase.initializeApp({
    apiKey: "AIzaSyCybtGa3ewvhu09JIzBCyodpbHJ-9infsI",
    authDomain: "lasm-go.firebaseapp.com",
    projectId: "lasm-go",
    storageBucket: "lasm-go.appspot.com",
    messagingSenderId: "69394339404",
    appId: "1:69394339404:web:40c6f58a252782a911c6e9",
});

const messaging = firebase.messaging();
if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/firebase-messaging-sw.js')
    .then(function(registration) {
        console.log('Service Worker Registered', registration);
        // Initialize Firebase Messaging with the registration object
        messaging.useServiceWorker(registration);
    })
    .catch(function(err) {
        console.log('Service Worker Registration Failed', err);
    });
}
