<script lang="ts">
    import q_yellow from "$lib/assets/images/general/questionmark_yellow.svg";
    import q_orange from "$lib/assets/images/general/questionmark_orange.svg";
    import q_pink from "$lib/assets/images/general/questionmark_pink.svg";
    import q_purple from "$lib/assets/images/general/questionmark_purple.svg";
    import { onMount } from "svelte";
    
    let qs = [q_yellow, q_orange, q_pink, q_purple];
    let bgcols = ["#fffffa", "#fffbfa", "#fffafe", "#fefaff"];

    let current_q_idx = $state(0);

    onMount(() => {
        const interval = setInterval(() => {
            current_q_idx = (current_q_idx + 1) % qs.length;
        }, 1000);

        return () => clearInterval(interval);
    });
    let { team } = $props()
    // export let team: any;

    let players: {name: string, email: string}[] = [];

    $effect(() => {
        players = [
            {name: team.player1name, email: team.player1email},
            {name: team.player2name, email: team.player2email},
            {name: team.player3name, email: team.player3email},
            {name: team.player4name, email: team.player4email},
            {name: team.player5name, email: team.player5email},
        ]
    });


    let oppened: boolean = $state(false);

    function close() {
        oppened = false;
    }

    function open() {
        oppened = true;
    }




</script>

<main on:click={open} style="--hover-color: {bgcols[current_q_idx]}">
    {#if oppened}
    <div class="popup-holder" on:click|stopPropagation={close}>
        <div class="popup" on:click|stopPropagation>
            <h2 class="name">{team.name}</h2>
            <ul class="player-list">
                {#each players as player, i}
                    {#if player.email || player.name}
                        <li class="player-item">
                            <div class="dot" style="animation-delay:-{players.length - i}s"></div>
                            <div class="text">
                                <span class="player-name">{player.name}</span>
                                <span class="player-email">{player.email}</span>
                            </div>
                        </li>
                    {/if}
                {/each}
            </ul>
        </div>
    </div>
    {/if}
    <div class="slider">
        <h2 class="name">{team.name}</h2>
        <div class="img-wrapper">
            <img src={qs[current_q_idx]} height="50px" />
        </div>
    </div>
</main>

<style lang="scss">
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
        
        .text {
            display: flex;
            flex-direction: row;
            justify-content: space-between;
            gap: 20px;
            flex-grow: 1;
        }
    }

    .player-item {
        display: flex;
        gap: 20px;
    }

    .player-name {
        font-weight: 600;
        /* width: 150px; */
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

</style>