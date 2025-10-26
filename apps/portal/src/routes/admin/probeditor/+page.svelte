<script lang="ts">
    import ProbBanner from "$lib/components/admin/ProbBanner.svelte";
    import type { CorrectorsResponse } from "$lib/pocketbase-types";
import { onMount } from "svelte";
    let { data } = $props();



    

</script>



<main>
    <div class="banners-scrollview">
        <div class="banners-holder">
            {#each data.probs as prob}
            <button class="banner-select">
                <ProbBanner prob={prob} user={data.user as CorrectorsResponse} />
            </button>
            {/each}
        </div>
    </div>

</main>

<style>
    main {
        height: 100%;
        padding: 0px;
    }
    .banners-scrollview {
        overflow-y: scroll;
        width: 280px;
        height: 100%;
        padding: 20px;
    }

    .banners-holder {
        display: flex;
        justify-content: flex-start;
        align-items: flex-start;
        flex-direction: column;
    }

    button.banner-select {
        all: unset;
        cursor: pointer;
    }




</style>