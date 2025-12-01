<script lang="ts">
  let { data } = $props();
  import { derived } from 'svelte/store';
  import type { ConstantsResponse } from '$lib/pocketbase-types';
  import { Latex } from 'shared';

  let constants: Record<string, ConstantsResponse> = $derived(Object.fromEntries(data.constants.map((c) => [c.id, c])));


  // initialData.map((c) => [c.id, { constant: c, edit: {} }])

  let groupedConstantIds = $derived(Object.fromEntries(Object.entries(
    Object.groupBy(
          Object.entries(constants),
          (item) => item[1].group ?? "Nezarazeno"
      )
    ).flatMap(([group, probs]) => probs === undefined ? [] : [[group, probs.map((e) => e[0])]])));

</script>
<main class="content">
  <div class="content-constants">
    <div class="content">
        {#each Object.entries(groupedConstantIds) as idGroup, i }
            <div class="table-wrapper" class:yellow={ i%2 == 0 } class:orange={ i%2 == 1 }>
                <h2 class="title">{idGroup[0]}</h2>
                <table>
                    <thead> 
                        <tr>
                            <td>Název</td>
                            <td>Symbol</td>
                            <td>Hodnota</td>
                            <td>Jednotka</td>
                        </tr>
                    </thead>
                    <tbody>
                        {#each idGroup[1] as constantId}
                            <tr>
                                <td class="name">
                                    <Latex text={constants[constantId].name} />
                                </td>
                                <td class="symbol">
                                    <Latex text={constants[constantId].symbol} />
                                </td>
                                <td class="value">
                                    <Latex text={constants[constantId].value.toString()} />
                                </td>
                                <td class="unit">
                                    <Latex text={constants[constantId].unit} />
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/each} 
    </div>
  </div>
</main>


<style lang="scss">
    main {
        display: flex;
        flex-direction: column;
        flex-grow: 1;
        // max-width: 750px;
        gap: 20px;
        box-sizing: border-box;
    }
    
    .content {
            flex-grow: 1;
            display: flex;
            flex-direction: column;
            min-height: 0px;
            overflow: hidden;
            padding: 100px;
            padding-top: 20px;
                    
        .content-constants {
            display: flex;
            flex-direction: column;
            gap: 20px;
            min-height: 0px;
            flex-grow: 1;

            .controls {
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: flex-start;
                gap: 15px;

                button {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    justify-content: center;
                    gap: 8px;
                    border-radius: 3px;
                    width: fit-content;
                    cursor: pointer;
                    padding: 5px 20px;
                    color: #333333;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;

                    i {
                        font-size: 20px;
                    }

                    p {
                        font-family: 'Fredoka';
                        font-size: 16px;
                        margin: 0px;
                    }
                }

                .add {
                    border: 3px var(--color-yellow) solid;
                    background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-yellow) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-yellow) 20%, black 80%);
                    }
                }

                .delete {
                    border: 3px var(--color-pink) solid;
                    background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-pink) 20%, black 80%);
                    }

                    &.inactive {
                        cursor: not-allowed;
                        color: #555555;
                        border-color: #AAAAAA;
                        background-color: color-mix(in srgb, #333333 5%, transparent 95%);
                    }
                    
                }

                .edit {
                    border: 3px var(--color-purple) solid;
                    background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);

                    &:hover {
                        background-color: color-mix(in srgb, var(--color-purple) 20%, transparent 80%);
                        color: color-mix(in srgb, var(--color-purple) 20%, black 80%);
                    }

                    &.inactive {
                        cursor: not-allowed;
                        color: #555555;
                        border-color: #AAAAAA;
                        background-color: color-mix(in srgb, #333333 5%, transparent 95%);
                    }
                    
                }
            }

            .content {
                overflow-y: auto;
                scrollbar-width: 0px;
                min-height: 0;
                display: flex;
                flex-direction: column;
                align-items: center;
                gap: 30px;
                flex-grow: 1;

                &::-webkit-scrollbar {
                    display: none;
                }
            
                .table-wrapper {
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    width: 100%;
                    background-color: #FCFCFC;

                    border: 3px lightgray solid;
                    border-radius: 5px;

                    padding: 10px;
                    padding-top: 0px;
                    
                    font-family: 'Lexend';

                    &.yellow {
                        background-color: color-mix(in srgb, var(--color-yellow) 3%, transparent 97%);
                        border-color: color-mix(in srgb, var(--color-yellow) 50%, white 50%);

                        tbody {
                            td {
                                background-color: #f0f0f0C0;
    
                            }
                        }
                    }

                    &.orange {
                        background-color: color-mix(in srgb, var(--color-orange) 3%, transparent 97%);
                        border-color: color-mix(in srgb, var(--color-orange) 70%, white 30%);

                        tbody {
                            td {
                                background-color: #f0f0f0C0;
    
                            }
                        }
                    }

                    .title {
                        margin: 20px;
                    }

                    table {
                        width: 100%;
                        
                        thead {
                            td {
                                text-align: center;
                                font-family: 'Fredoka';
                                font-weight: 600;
                                font-size: 18px;
                            }
                        }
    
                        tbody {
                            tr {
                                
                                td {
                                    padding: 3px 10px;
                                    border-radius: 3px;
                                    cursor: pointer;
                                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.15s;
                                    
                                    &.name {
                                        border-left: 3px solid transparent;
                                    }
                                }
        
                                &:hover {
                                    td {
                                        background-color: color-mix(in srgb, var(--color-pink) 20%, transparent 80%);
                                    }
                                }
    
                                &.selected {
                                    td {
                                        background-color: color-mix(in srgb, var(--color-pink) 10%, transparent 90%);
                                        
                                        &.name {
                                            border-left: 3px solid var(--color-pink);    
                                        }
                                    }
                                }
                            }
                        }

                    }


                }
            }

            .edit-dialog-wrapper {
                display: none;
                position: fixed;
                top: 0px;
                left: 0px;
                width: 100%;
                height: 100%;
                
                background-color: #00000000;
                backdrop-filter: blur(0.0px);
                z-index: -100;
                transition: background-color cubic-bezier(0.215, 0.610, 0.355, 1) 0.5s;
                    
                &.active {
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    background-color: #00000033;
                    backdrop-filter: blur(2.0px);
                    z-index: 100;
                    
                }
            }
        }

    }
</style>
