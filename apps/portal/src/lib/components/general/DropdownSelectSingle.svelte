<script lang="ts">
    import type { DropdownItem } from '$lib/types';

    type OnChangeFunction = (id: string) => void;

    let {
        selected,
        possible,
        onchangeFunction
    }: {
        selected: DropdownItem,
        possible: DropdownItem[],
        onchangeFunction: OnChangeFunction
    } = $props();

    let opened = $state(false);

    let dropdownHeight = $state(0);
</script>

<main>
    <button class="dropdown" onclick={() => { opened = !opened; }} class:open={ opened }>
        <p>{selected.value}</p>
        <i class="fa-solid fa-angle-down"></i>
    </button>
    <div class="content" class:open={ opened } style="--wrapper-height: { dropdownHeight }px">
        <div class="wrapper" bind:clientHeight={ dropdownHeight }>
            {#each possible as possibility }
                <button class:selected={ selected.id === possibility.id } onclick={() => { onchangeFunction(possibility.id); opened = false; }}>
                    {possibility.value}
                </button>
            {/each}
        </div>
    </div>
</main>


<style lang="scss">
    main {
        display: flex;
        flex-direction: column;
        gap: 8px;
        align-items: flex-end;

        button.dropdown {
            all: unset;
            cursor: pointer;
            padding: 5px 15px;
            background-color: #F5F5F5;
            border-radius: 3px;
            color: #333333;
            display: flex;
            flex-direction: row;
            justify-content: space-between;
            align-items: center;
            gap: 25px;
            width: fit-content;
            
            p {
                font-family: 'Fredoka';
                font-size: 20px;
                font-weight: 500;
                margin: 0px;
            }
            
            i {
                transform: rotate(180deg);
                transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
            }

            &.open {

                i {
                    transform: rotate(0deg);
                }
            }
        }

        .content {
            min-height: 0px;
            height: 0px;
            transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
            overflow: hidden;
            width: 100%;
            
            
            &.open {
                height: var(--wrapper-height);
            }
            
            .wrapper {
                display: flex;
                flex-direction: column;
                padding: 10px 5px;
                align-items: flex-start;
                gap: 10px;
                box-sizing: border-box;
                width: 100%;
                background-color: #FAFAFA;
                border-radius: 5px;
                // padding-top: 10px;
                // padding-bottom: 10px;

                button {
                    all: unset;
                    cursor: pointer;
                    font-family: 'Fredoka';
                    width: 100%;
                    box-sizing: border-box;
                    padding: 3px 15px;
                    border: 3px transparent solid;
                    border-radius: 3px;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.2s;

                    &.selected {
                        background-color: color-mix(in srgb, var(--color-pink) 10%, transparent 90%);
                        border-color: var(--color-pink);
                    }

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-pink) 10%, transparent 90%);
                    }
                }
            }
        }
    }
</style>