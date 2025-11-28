<script lang="ts">
    import { onMount } from "svelte";

    let dotCount = $state(1);

    onMount(() => {
        setTimeout(() => {
            setInterval(() => {
                if (dotCount == 6) {
                    dotCount = 1;
                } else {
                    dotCount += 1;
                }
            }, 300);
        }, 1000);
    });


</script>

<main>
    <h1>Spojení ztraceno{Array(dotCount).fill(".").join("")}</h1>
    <p>Vyčkejte chvíly a obnovte stránku (F5)</p>
</main>

<style lang="scss">
    main {
        background-color: black;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        
        h1 {
            font-family: monospace;
            color: white;
            text-align: center;
            text-wrap: nowrap;
            opacity: 0;
            animation: emerge 1s cubic-bezier(0.215, 0.610, 0.355, 1) forwards 0.5s;
        }

        p {
            font-family: monospace;
            color: white;
            text-wrap: nowrap;
            text-align: center;
            font-size: 16px;
            opacity: 0;
            animation: emerge 1s cubic-bezier(0.215, 0.610, 0.355, 1) forwards 1.2s;
        }

        @keyframes emerge {
            0% {
                transform: translateY(-30px);
                opacity: 0;
            }
            100% {
                transform: translateY(0px);
                opacity: 1;
            }
        }
    }
</style>