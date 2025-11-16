<script lang="ts">
    let { color }: { color: string } = $props();

    let size = $state(0);
</script>

<main style="--main-color: {color}; --size: {size}px">
    <div class="block one"></div>
    <div class="block two"></div>
</main>

<style lang="scss">
    main {
        width: 100%;
        height: 100%;
        padding-bottom: 100%;
        box-sizing: border-box;
        position: relative;

        .block {
            background-color: var(--main-color);
            // background-color: red;
            border-radius: 15%;
            position: absolute;
            animation: animation 1s cubic-bezier(0.215, 0.610, 0.355, 1) infinite;
            width: 40%;
            height: 40%;
            
            &.one {
                animation-delay: 0s;
            }

            &.one {
                animation-delay: -0.5s;
            }
        }
    }

    @keyframes animation {
        0% {
            top: 0%;
            left: 0%;
        }
        25% {
            top: 0%;
            left: 60%
        }
        50% {
            top: 60%;
            left: 60%
        }
        75% {
            top: 60%;
            left: 0%;
        }
        100% {
            top: 0%;
            left: 0%;
        }
    }
</style>
