<script lang="ts">
    export let hovered: boolean = false;
    export let text: string = "hoverable";


</script>

<p class:hovered = {hovered}>
    {#each text as char, i}
        <span class="char" style="animation-delay: {i % 2 == 0 ? "0s" : "-0.5s"};">{char}</span>
    {/each}
</p>


<style>
    .char {
        color: white;
        font-family: "Lexend";
    }
    p.hovered .char {
        animation: hover-anim 1s steps(1) infinite;
        font-family: "Lexend";
    }

    p {
        margin: 0px;
    }

    @keyframes hover-anim {
        0% {
            color: #EBAD00;
        }
        50% {
            color: #EB6E00;
        }
    }
</style>