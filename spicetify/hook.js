(function MoodThemeClient() {
    if (!Spicetify?.Player) {
        setTimeout(NexusClient, 1000);
        return;
    }

    console.log("Success!: Extension injected, listening to songs...");

    Spicetify.Player.addEventListener("songchange", () => {
        const data = Spicetify.Player.data;
        
        if (!data?.item) return;

        const songName = data.item.name;
        console.log(`Playing: ${songName}. notifiying to the backend...`);

        fetch("http://localhost:8080/theme", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ song: songName })
        })
        .then(async (response) => {
            const data = await response.json();
            
            if (response.ok) {
                console.log(`Success or wait!: ${data.message}`);
            } else {
                console.warn(`Backend warning: ${data.message || data.error}`);
            }
        })
        .catch(error => {
            console.error("The server has crashed or is off", error);
        });
    });
})();