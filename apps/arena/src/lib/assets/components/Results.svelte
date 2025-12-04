<script lang="ts">
    import { currentState } from "$lib/stores/state";
</script>

<main>
    <h1 class="name">{$currentState.teamName}</h1>
    <div class="stats">
        <div class="rank property-wrapper">
            <h2 class="value">{$currentState.rank}.</h2>
            <p class="name">Umístění</p>
        </div>
        <div class="money property-wrapper">
            <h2 class="value">{$currentState.money}</h2>
            <p class="name">Body</p>
        </div>
    </div>
    {#if $currentState.rank > 15}
        <h2 class="message">Pokud jste se umístili v top 15, dostanete potvrzující email s dalšími instrukcemi o postupu do finálového kola</h2>
    {:else}
        <h2 class="message">Vypadá to, že jste se umístili v top 15! Očekávejte tedy potvrzující email s dalšími instrukcemi o postupu do finálového kola</h2>
    {/if}
</main>

<style lang="scss">
    main {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        width: 100%;

        h1.name {
            font-family: 'Lexend';
            font-size: 45px;
            color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
        }

        .stats {
            width: 100%;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: space-evenly;

            .property-wrapper {
                display: flex;
                flex-direction: column;
                gap: 15px;
                justify-content: center;
                align-items: center;
                padding-top: 100px;

                .value {
                    font-family: 'Lexend';
                    font-size: 60px;
                    margin: 0px;
                    padding: 10px 60px;

                    border-bottom: 3px color-mix(in srgb, var(--color-purple) 30%, black 70%) dashed;
                }

                .name {
                    font-family: 'Fredoka';
                    font-size: 25px;
                    color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
                    font-weight: 500;
                    margin: 0px;
                }
            }



            .money {
                .value {
                    color: var(--color-pink);
                }
            }
            .rank {
                .value {
                    color: var(--color-yellow);
                }
            }
        }

        .message {
            padding-top: 50px;
            font-family: 'Fredoka';
            font-weight: 500;
            color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
        }
    }
</style>