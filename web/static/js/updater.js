// Auto-refresh the page every 2 seconds to get updated status
const REFRESH_INTERVAL = 2000; // milliseconds

function refreshPage() {
    // Reload the page to get fresh data from the server
    window.location.reload();
}

// Set up auto-refresh
setInterval(refreshPage, REFRESH_INTERVAL);

// Optional: Add visibility API to pause updates when tab is not visible
document.addEventListener('visibilitychange', function() {
    if (document.hidden) {
        console.log('Page hidden - updates will continue');
    } else {
        console.log('Page visible - refreshing now');
        refreshPage();
    }
});
