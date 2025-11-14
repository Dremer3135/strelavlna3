<script lang="ts">
    let { children } = $props();

    let hovered: boolean = $state(false);


</script>

<main onmouseenter={() => { hovered = true; }} onmouseleave={() => { hovered = false; }}>
    <!-- <div class="wrapper"> -->
    {@render children()}
    <div class="cube"></div>
    <!-- </div> -->
</main>

<style lang="scss">
    main {
        position: relative;
        width: fit-content;
        height: fit-content;

        .cube {
            position: absolute;
            width: 10px;
            height: 10px;
            border-radius: 3px;
        }
    }
</style>