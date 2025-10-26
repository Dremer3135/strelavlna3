<script lang="ts">
    import type { ProbsRecord } from "$lib/pocketbase-types";
    import type { AppUser } from "../../../app";
    let { prob, user }: {prob: ProbsRecord, user: AppUser} = $props();


    // 0 - free
    // 1 - someone has it
    // 2 current user has it
    let authorType = $derived(prob.author == "" ? 0 : prob.author == user.id ? 2 : 1);






</script>

<main class:free={authorType == 0} class:owned={authorType == 1} class:mine={authorType == 2}>
    <h3 class="name">{prob.name}</h3>
    <p class="diff">[{prob.diff}]</p>
</main>

<style>
    main {
        background-color: white;
        border: 3px white solid;
        /* border-width: 2px; */
        border-radius: 3px; 
        font-family: 'Lexend';
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 5px 20px;
        width: 200px;
        transition: all 0.3s cubic-bezier(0.215, 0.610, 0.355, 1);
    }

    main.free {
        border-color: gray;
        background-color: color-mix(in srgb, gray 5%, transparent 95%);
    }
    main.free:hover {
        background-color: color-mix(in srgb, gray 20%, transparent 80%);
        transition: all 0s;
    }

    main.owned {
        border-color: black;
        background-color: color-mix(in srgb, gray 5%, transparent 95%);
    }
    main.owned:hover {
        background-color: color-mix(in srgb, gray 20%, transparent 80%);
        transition: all 0s;
    }

    main.mine {
        border-color: var(--color-pink);
        background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
    }
    main.mine:hover {
        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
        transition: all 0s;
    }

    h3 {
        margin: 0px;
        text-wrap: wrap;
    }
    p {
        margin: 0px;
    }

</style>