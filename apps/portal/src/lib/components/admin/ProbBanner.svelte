<script lang="ts">
    import type { ProbsRecord } from "$lib/pocketbase-types";
    import type { AppUser } from "../../../app";
    let { prob, user, selected }: {prob: ProbsRecord, user: AppUser, selected: boolean} = $props();


    // 0 - free
    // 1 - someone has it
    // 2 current user has it
    let authorType = $derived(prob.author == "" ? 0 : prob.author == user.id ? 2 : 1);








</script>

<main class:free={authorType == 0} class:owned={authorType == 1} class:mine={authorType == 2} class:selected={selected}>
    <div class="background"></div>
    <h3 class="name">{prob.name}</h3>
    <p class="diff">[{prob.diff}]</p>
</main>

<style lang="scss">
    main {
        position:relative;
        /* background-color: white; */
        border: 3px white solid;
        border-radius: 3px; 
        font-family: 'Lexend';
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 5px 20px;
        width: 100%;
        max-width: 300px;
        transition: all 0.3s cubic-bezier(0.215, 0.610, 0.355, 1);
        overflow: hidden;
    }

    .background {
        top: 0px;
        left: 0px;
        width: 100%;
        height: 100%;
        box-sizing: border-box;
        position: absolute;
        background-color: white;
        z-index: -1;
        border-right: 10px black solid;
        transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.2s;
        transform: translateX(10px);
        will-change: transform;
        
    }

    .selected .background {
        transform: translateX(0px);
    }

    main.free {
        border-color: lightgray;
        background-color: color-mix(in srgb, lightgray 5%, transparent 95%);

        .background {
            border-color: lightgray;
        }
    }
    main.free:hover {
        background-color: color-mix(in srgb, lightgray 20%, transparent 80%);
        /* transition: all 0s;   */
    }

    main.owned {
        border-color: black;
        background-color: color-mix(in srgb, gray 5%, transparent 95%);
        
        .background {
            border-color: black;
        }
    }
    main.owned:hover {
        background-color: color-mix(in srgb, gray 20%, transparent 80%);
        /* transition: all 0s; */
    }

    main.mine {
        border-color: var(--color-orange);
        background-color: color-mix(in srgb, var(--color-orange) 5%, transparent 95%);

        .background {
            border-color: var(--color-orange);
        }
    }
    main.mine:hover {
        background-color: color-mix(in srgb, var(--color-orange) 20%, transparent 80%);
        /* transition: all 0s; */
    }

    h3 {
        margin: 0px;
        text-wrap: wrap;
    }
    p {
        margin: 0px;
        transform: translateY(-2px);
    }

    .diff {
        font-size: 20px;
    }

</style>