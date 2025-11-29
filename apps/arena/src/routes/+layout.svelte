<script lang="ts">
    import Navbar from '$lib/assets/components/Navbar.svelte';
	import favicon from '$lib/assets/favicon.png';
    import { wsConnected } from '$lib/stores/state';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<main>
    <!-- {#if $wsConnected}
    <Navbar />
    {/if} -->
    {@render children?.()}
</main>


<style>
    @import url('https://fonts.googleapis.com/css2?family=Lexend:wght@100..900&display=swap');
    @font-face {
        font-family: "Fredoka";
        
        src: url('/fonts/Fredoka-VF.ttf') format('truetype');

        font-style: normal;
        /* font-weight: normal; */
    }

    :root {
        --color-purple: #9500EB;
        --color-pink: #EB0072;
        --color-orange: #EB6E00;
        --color-yellow: #EBAD00;
        --color-blue: #3c008b;
        --color-lightblue:#003F88;
        --color-gray: #353538;
        --color-light-gray: #f8f7ff;
    }

    main {
        display: flex;
        flex-direction: column;
        min-height: 0px;
        width: 100vw;
        height: 100vh;
    }

    main::-webkit-scrollbar {
        display: none;
    }
</style>