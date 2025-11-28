<script lang="ts">
    import logo from "$lib/assets/images/logo.svg"
    import Digit from "./Digit.svelte";
    import { currentState } from "$lib/stores/state";
    import { onMount } from "svelte";

    let remaining = $state(0);
    
    let hours: string = $derived(Math.floor(remaining / 1000 / 60 / 60).toString().padStart(2, "0"));
    let minutes: string = $derived((Math.floor(remaining / 1000 / 60) % 60).toString().padStart(2, "0"));
    let seconds: string = $derived((Math.floor(remaining / 1000) % 60).toString().padStart(2, "0"));
    
    
    function updateClock() {
        if (Date.now() > $currentState.start.getTime()) {
            remaining = Math.abs($currentState.end.getTime() - Date.now());
        } else {
            remaining = Math.abs($currentState.start.getTime() - Date.now());
        }
    }

    onMount(() => {
        setInterval(updateClock, 1000);
    });


</script>

<main>
    <div class="logo">
        <img src={logo} alt="logo">
    </div>
    <div class="timer">
        <h3>{hours}:{minutes}:{seconds}</h3>
    </div>
</main>

<style lang="scss">
    main {
        width: 100%;
        padding: 15px 50px;
        box-sizing: border-box;
        background: linear-gradient(to right, var(--color-lightblue) 0%, var(--color-pink) 150%);
        display: flex;
        align-items: center;
        justify-content: space-between;

        .logo {
            img {
                height: 50px;
            }
        }

        .timer {
            width: 110px;
            text-align: center;
            display: flex;
            flex-direction: row;
            align-items: center;
            gap: 8px;
            height: 50px;

            h3 {
                font-family: 'Lexend';
                font-size: 23px;
                font-weight: 700;
                color: white;
                margin: 0px;
            }
        }
    }
</style>