<script lang="ts">


    let { value = $bindable(), groups }: { value: string | undefined, groups: string[] } = $props();

    let suggestions: string[] = $derived(groups.filter(item => item.toLowerCase().includes(value?.toLowerCase() ?? "")));

    let isFocused: boolean = $state(false);

    let suggestionsHeight: number = $state(0);

</script>


<main>
    <input type="text" bind:value={value} onfocus={() => { isFocused = true; }} onblur={() => { setTimeout(() => {isFocused = false; }, 100); }} placeholder="Group">
    <!-- {#if isFocused} -->
    <div class="suggestions" class:collapsed={!isFocused} style="--wrapper-height: { suggestionsHeight }px;">
        <div class="wrapper" bind:clientHeight={ suggestionsHeight }>
            {#each suggestions as suggestion}
                <button class="suggestion" onclick={() => {value = suggestion; isFocused = false; }}>
                    <div class="indicator"></div>
                    {suggestion}
                </button>
            {/each}
        </div>
    </div>
    <!-- {/if} -->
</main>

<style lang="scss">
    main {
        input {
            all: unset;
            cursor: text;
            font-family: 'Fredoka';
            border: 3px lightgray solid;
            padding: 4px 12px;
            border-radius: 3px;
            background-color: #FAFAFA;
            width: 450px;

            transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

            &:focus {
                border-color: var(--color-yellow);
                background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);
            }
        }

        .suggestions {
            min-height: 0;
            height: var(--wrapper-height);
            transition: height cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
            overflow: hidden;
            &.collapsed {
                height: 0px;
            }

            .wrapper {
                padding: 20px;
                display: flex;
                flex-direction: column;
                justify-content: center;
                align-items: flex-start;
                gap: 10px;

                button {
                    all: unset;
                    cursor: pointer;
                    font-family: 'Fredoka';
                    font-size: 17px;
                    font-weight: 600;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    gap: 10px;
                    color: #666666;

                    .indicator {
                        width: 6px;
                        height: 6px;
                        border-radius: 50%;
                        background-color: #AAAAAA;
                    }

                    &:hover {
                        color: black;

                        .indicator {
                            width: 8px;
                            height: 8px;
                            transform: translateX(-1px);
                            background-color: color-mix(in srgb, var(--color-pink) 80%, white 20%);
                        }
                    }
                }
            }
        }
    }
</style>