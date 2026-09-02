<script lang="ts">
    import q_yellow from "$lib/assets/images/register/questionmark_yellow.svg?url";
    import q_orange from "$lib/assets/images/register/questionmark_orange.svg?url";
    import q_pink from "$lib/assets/images/register/questionmark_pink.svg?url";
    import q_purple from "$lib/assets/images/register/questionmark_purple.svg?url";
    import { onMount } from "svelte";
    import { pocketbase } from "$lib/pocketbase";
    import { teamsStore } from "$lib/stores/register";
    import { myTeamsStore } from "$lib/stores/myTeams";
    
    let qs = [q_yellow, q_orange, q_pink, q_purple];
    let bgcols = ["#fffffa", "#fffbfa", "#fffafe", "#fefaff"];

    let current_q_idx = $state(0);

    onMount(() => {
        const interval = setInterval(() => {
            current_q_idx = (current_q_idx + 1) % qs.length;
        }, 1000);

        return () => clearInterval(interval);
    });

    let { team }: { team: any } = $props();

    let displayTeam = $state({ ...team });
    $effect(() => {
        displayTeam = { ...team };
    });

    let players = $derived([
        {name: displayTeam.player1name, email: displayTeam.player1email},
        {name: displayTeam.player2name, email: displayTeam.player2email},
        {name: displayTeam.player3name, email: displayTeam.player3email},
        {name: displayTeam.player4name, email: displayTeam.player4email},
        {name: displayTeam.player5name, email: displayTeam.player5email},
    ]);

    let oppened: boolean = $state(false);
    let isEditing: boolean = $state(false);
    let isSaving: boolean = $state(false);
    let error: string = $state("");

    let name: string = $state("");
    let editablePlayers = $state<{ name: string; email: string }[]>([]);

    function startEdit() {
        name = displayTeam.name ?? "";
        editablePlayers = [
            { name: displayTeam.player1name ?? "", email: displayTeam.player1email ?? "" },
            { name: displayTeam.player2name ?? "", email: displayTeam.player2email ?? "" },
            { name: displayTeam.player3name ?? "", email: displayTeam.player3email ?? "" },
            { name: displayTeam.player4name ?? "", email: displayTeam.player4email ?? "" },
            { name: displayTeam.player5name ?? "", email: displayTeam.player5email ?? "" },
        ];
        error = "";
        isEditing = true;
    }

    async function save() {
        if (!name.trim()) {
            error = "Jméno týmu je potřeba.";
            return;
        }
        isSaving = true;
        error = "";

        const updated = {
            name: name.trim(),
            player1name: editablePlayers[0].name.trim(),
            player1email: editablePlayers[0].email.trim(),
            player2name: editablePlayers[1].name.trim(),
            player2email: editablePlayers[1].email.trim(),
            player3name: editablePlayers[2].name.trim(),
            player3email: editablePlayers[2].email.trim(),
            player4name: editablePlayers[3].name.trim(),
            player4email: editablePlayers[3].email.trim(),
            player5name: editablePlayers[4].name.trim(),
            player5email: editablePlayers[4].email.trim(),
        };

        try {
            await pocketbase.collection("teams").update(team.id, updated);
            
            Object.assign(displayTeam, updated);
            Object.assign(team, updated);
            teamsStore.update(teams => teams.map(t => t.id === team.id ? { ...t, ...updated } : t));
            myTeamsStore.update(teams => teams.map(t => t.id === team.id ? { ...t, ...updated } : t));

            isEditing = false;
        } catch (err) {
            console.error(err);
            error = "Nepodařilo se uložit změny.";
        } finally {
            isSaving = false;
        }
    }

    function close() {
        oppened = false;
        isEditing = false;
        error = "";
    }

    function open() {
        oppened = true;
    }
</script>

<main on:click={open} style="--hover-color: {bgcols[current_q_idx]}">
    {#if oppened}
    <div class="popup-holder" on:click|stopPropagation={close}>
        <div class="popup" on:click|stopPropagation>
            {#if !isEditing}
                <h2 class="name">{displayTeam.name}</h2>
                <ul class="player-list">
                    {#each players as player, i}
                        {#if player.email || player.name}<li class="player-item"><div class="dot" style="animation-delay:-{players.length - i}s"></div><span class="player-name">{player.name}</span><span class="player-email">{player.email}</span></li>{/if}
                    {/each}
                </ul>
                <button class="edit-btn" on:click={startEdit}>Upravit</button>
            {:else}
                <form on:submit|preventDefault={save}>
                    <input class="name-input" type="text" placeholder="Název týmu" bind:value={name} required />
                    {#each editablePlayers as player, i}
                        <div class="edit-row">
                            <input type="text" placeholder={`Jméno hráče ${i + 1}`} bind:value={player.name} />
                            <input type="email" placeholder={`Email hráče ${i + 1}`} bind:value={player.email} />
                        </div>
                    {/each}
                    {#if error}<p class="error">{error}</p>{/if}
                    <div class="buttons">
                        <button type="button" on:click={() => (isEditing = false)} disabled={isSaving}>Zrušit</button>
                        <button type="submit" disabled={isSaving}>{isSaving ? "Ukládání..." : "Uložit"}</button>
                    </div>
                </form>
            {/if}
        </div>
    </div>
    {/if}
    <div class="slider">
        <h2 class="name">{displayTeam.name}</h2>
        <div class="img-wrapper">
            <img src={qs[current_q_idx]} height="50px" />
        </div>
    </div>
</main>

<style>
    main {
        background-color: #f8f8f8;
        padding: 0px 50px;
        cursor: pointer;
        height: 60px;
        overflow: hidden;
        transition: background-color 0.3s ease-out;
    }
    .slider {
        transform: translateY(0px);
        transition: transform 0.3s cubic-bezier(0.33, 1, 0.68, 1);
    }
    main:hover .slider {
        transform: translateY(-60px);
    }

    main:hover {
        background-color: var(--hover-color);
    }

    .img-wrapper {
        height: 60px;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .popup-holder {
        all: unset;
        position: fixed;
        top: 0;
        left: 0;
        height: 100%;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 100;

        background-color: #00000033;
        pointer-events: auto;

        cursor: auto;
    }
    .popup {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;

        background-color: white;

        border-radius: 10px;

        padding: 20px 80px; 
    }

    .name {
        font-family: "Fredoka";
        font-weight: 600;
        font-size: 18px;
        color: #28476d;
        flex-shrink: 1;
        height: 60px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0px;
    }

    .popup .name {
        font-family: "Fredoka";
        font-weight: 600;
        font-size: 40px;
        color: #002c5e;
        margin-top: 0px;
        margin-bottom: 35px;
    }

    .player-list {
        margin: 0px;
        gap: 20px;
        list-style: none;
        padding: 0px;
    }

    .player-list li {
        font-family: "Lexend";
        font-size: 18px;
        margin-bottom: 10px;
        display: flex;
        align-items: center;
        gap: 10px;

    }

    .player-item {
        display: flex;
        gap: 20px;
    }

    .player-name {
        font-weight: 600;
        width: 150px;
    }

    .player-email {
        color: #555;
    }
    
    .dot {
        width: 11px;
        height: 11px;
        border-radius: 3px;
        transform: translateY(1px);
        animation: list-animation 4s infinite steps(1, end);
    }


    @keyframes list-animation {
        0% {
            background-color: #EBAD00;
        }
        25% {
            background-color: #EB6E00;
        }
        50% {
            background-color: #EB0072;
        }
        75% {
            background-color: #9500EB;
        }
    }

    .edit-btn {
        font-family: "Lexend";
        font-size: 14px;
        margin-top: 20px;
        padding: 8px 20px;
        border-radius: 0px;
        border: 5px solid #d80000;
        background: #f8f8f8;
        cursor: pointer;
        color: #002c5e;
        font-weight: 500;
    }
    .edit-btn:hover {
        background: #002c5e;
        color: white;
    }
    form {
        display: flex;
        flex-direction: column;
        gap: 10px;
        width: 100%;
    }
    .name-input {
        font-family: "Fredoka";
        font-size: 24px;
        font-weight: 600;
        color: #002c5e;
        border: none;
        border-bottom: 2px solid #002c5e;
        padding: 5px;
        margin-bottom: 10px;
        outline: none;
    }
    .edit-row {
        display: flex;
        gap: 10px;
    }
    .edit-row input {
        font-family: "Lexend";
        padding: 8px 10px;
        border: 1px solid #ddd;
        border-radius: 6px;
        flex: 1;
        outline: none;
        width: 100%;
        box-sizing: border-box;
    }
    .edit-row input:focus {
        border-color: #002c5e;
    }
    .buttons {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 10px;
    }
    .buttons button {
        font-family: "Lexend";
        padding: 8px 18px;
        border-radius: 6px;
        border: 1px solid #ddd;
        background: #f8f8f8;
        cursor: pointer;
    }
    .buttons button[type="submit"] {
        background: #002c5e;
        color: white;
        border: none;
    }
    .error {
        color: red;
        font-family: "Lexend";
        font-size: 14px;
        margin: 0;
    }

</style>