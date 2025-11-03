<script lang="ts">

</script>

<main>
    <div class="content">
        <div class="first"></div>
        <div class="second"></div>
    </div>
</main>

<style lang="scss">
    main {
        position: relative;
        padding: 20px;
        // background-color: #FAFAFA;
        max-width: fit-content;
        max-height: fit-content;
        border-radius: 5px;

        
        .content {
            position: relative;
            width: 60px;
            height: 60px;
        
            div {
                position: absolute;
                width: 25px;
                height: 25px;
                border-radius: 5px;
            }


            .first {
                animation: loading 1s cubic-bezier(0.215, 0.610, 0.355, 1) infinite;
                background-color: #007bff; // Primary color
            }

            .second {
                animation: loading 1s cubic-bezier(0.215, 0.610, 0.355, 1) infinite -0.5s; // -1s delay offsets the start
                background-color: #ffc107; // A contrasting color
            }
        }
    }




    @keyframes loading {
        0% {
            top: 0px;
            left: 0px;
            background-color: var(--color-yellow);
            box-shadow: 0px 0px 10px 5px color-mix(in srgb, var(--color-yellow) 30%, transparent 70%);
        }
        25% {
            top: 0px;
            left: 35px;
            background-color: var(--color-purple);
            box-shadow: 0px 0px 10px 5px color-mix(in srgb, var(--color-purple) 30%, transparent 70%);
        }
        50% {
            top: 35px;
            left: 35px;
            background-color: var(--color-orange);
            box-shadow: 0px 0px 10px 5px color-mix(in srgb, var(--color-orange) 30%, transparent 70%);
        }
        75% {
            top: 35px;
            left: 0px;
            background-color: var(--color-pink);
            box-shadow: 0px 0px 10px 5px color-mix(in srgb, var(--color-pink) 30%, transparent 70%);
        }
        100% {
            top: 0px;
            left: 0px;
            background-color: var(--color-yellow);
            box-shadow: 0px 0px 10px 5px color-mix(in srgb, var(--color-yellow) 30%, transparent 70%);
        }
    }




</style>