<script lang="ts">
	import Button from '$lib/assets/components/Button.svelte';
	import ProbSelect from '$lib/assets/components/ProbSelect.svelte';
	import { probs } from '$lib/stores/probs';
	// import { onMount } from "svelte"; // No longer needed for this
	import type { MessageType } from '$lib/types';
	import { focusedProb, currentState } from '$lib/stores/state';
	import { ProbContent } from 'shared';
	import AnswerInput from '$lib/assets/components/AnswerInput.svelte';
	import Chat from '$lib/assets/components/Chat.svelte';
	import { getProbLastAnswer, isProbSolved, isProbUngraded } from '$lib/utils';
    import { untrack } from 'svelte';

	let {
		buy,
		sell,
		focus,
		chat
	}: {
		buy: (diff: string) => void;
		sell: (id: string) => void;
		focus: (id: string) => void;
		chat: (probId: string, message: Omit<MessageType, 'sentTime'>) => void;
	} = $props();

	let isWindowFocused = $state(true);

	function handleVisibilityChange() {
		isWindowFocused = !document.hidden;
	}

	function handleWindowBlur() {
		isWindowFocused = false;
	}

	function handleWindowFocus() {
		if (!document.hidden) {
			isWindowFocused = true;
		}
	}

	const difficulties = ['A', 'B', 'C'];

	let buySelectOpened: boolean = $state(false);

	$effect(() => {
        let fProb = untrack(() => $focusedProb);
        if (fProb) {
            if (isWindowFocused) {
                chat(fProb.id, {
                    origin: "sent",
                    type: "window-focus",
                    value: "blur " + untrack(() => $currentState.myId)
                });
            } else {
                chat(fProb.id, {
                    origin: "sent",
                    type: "window-focus",
                    value: "focus " + untrack(() => $currentState.myId)
                });
            }
        }
	});

	let buyHovered = $state(Array(difficulties.length).fill(false));
	let sellHovered = $state(false);
</script>

<!-- Combines both methods for the most robust detection -->
<svelte:document on:visibilitychange={handleVisibilityChange} />
<svelte:window on:blur={handleWindowBlur} on:focus={handleWindowFocus} />

<main>
	<div class="left-panel">
		<div class="team-stats">
			<h2 class="name">{$currentState.teamName}</h2>
			<h3 class="money">{$currentState.money} DC</h3>
		</div>
		<div class="prob-selector">
			{#each Object.values($probs).filter((prob) => !isProbSolved(prob)) as prob}
				<ProbSelect
					prob={prob}
					onSelect={() => {
						focus(prob.id);
					}}
				/>
			{/each}
            <div class="separator"></div>
			{#each Object.values($probs).filter((prob) => isProbSolved(prob)) as prob}
				<ProbSelect
					prob={prob}
					onSelect={() => {
						focus(prob.id);
					}}
				/>
			{/each}
		</div>

		<div class="controls">
			<div class="buy">
				{#each difficulties as diff, i}
					<Button
						disabled={$currentState.pricesBuy[i] > $currentState.money ||
							$currentState.probsRemaining[i] == 0}
						theme={['yellow', 'orange', 'pink', 'purple'][i % 4] as
							| 'yellow'
							| 'orange'
							| 'pink'
							| 'purple'}
						onclick={() => {
							if ($currentState.pricesBuy[i] > $currentState.money || $currentState.probsRemaining[i] == 0)
								return;
							buySelectOpened = !buySelectOpened;
							buy(diff);
						}}
						onmouseenter={() => {
							buyHovered[i] = true;
						}}
						onmouseleave={() => {
							buyHovered[i] = false;
						}}
					>
						<p>Koupit [{diff}]</p>
						<span class="anchor">
							<p class="tooltip" class:visible={buyHovered[i]}>
								Zbývá: {$currentState.probsRemaining[i] == -1 ? '∞' : $currentState.probsRemaining[i]}
							</p>
						</span>
					</Button>
				{/each}
			</div>
			<div class="sell">
				<Button
					disabled={!$focusedProb}
					theme="pink"
					onmouseenter={() => {
						sellHovered = true;
					}}
					onmouseleave={() => {
						sellHovered = false;
					}}
					onclick={() => {
						if ($focusedProb) {
							sell($focusedProb.id);
						}
					}}
				>
					<i class="fa-solid fa-trash-can" />
					<p>Prodat</p>
					<span class="anchor">
						<p class="tooltip" class:visible={$focusedProb && sellHovered}>
							Prodat úlohu <span class="bold">{$focusedProb?.name}</span>
						</p>
					</span>
				</Button>
			</div>
		</div>
	</div>

	<div class="content">
		{#if $focusedProb}
			<div class="prob">
				<ProbContent
					content={$focusedProb}
					onCopy={(text: string) => {
						console.log('copying');

						if (!$focusedProb) return;

						chat($focusedProb.id, {
							type: 'copy',
							value: text,
							origin: 'sent'
						});
					}}
				/>
				<AnswerInput
					disabled={isProbUngraded($focusedProb) || isProbSolved($focusedProb)}
					placeholder={getProbLastAnswer($focusedProb) !== undefined
						? (getProbLastAnswer($focusedProb) as string)
						: 'Odpověď'}
					submitAnswer={(answer: string) => {
						if (!$focusedProb) return;

						chat($focusedProb.id, {
							type: 'answer',
							value: answer,
							origin: 'sent'
						});
					}}
					onPaste={(text: string) => {
						if (!$focusedProb) return;

						chat($focusedProb.id, {
							type: 'paste',
							value: text,
							origin: 'sent'
						});
					}}
				/>
			</div>
			<div class="chat">
				<Chat
                    type="player"
					prob={$focusedProb}
					send={(message) => {
						if (!$focusedProb) return;

						chat($focusedProb.id, message);
					}}
				/>
			</div>
		{/if}
	</div>
</main>

<style lang="scss">

    main {
        flex-grow: 1;
        display: flex;
        flex-direction: row;
        min-height: 0px;

        .left-panel {
            width: 350px;
            display: flex;
            flex-direction: column;
            min-height: 0px;


            .team-stats {
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: space-between;
                padding: 20px;
                box-sizing: border-box;
                background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                border-right: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                border-bottom: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                gap: 10px;

                .name {
                    font-family: 'Lexend';
                    font-weight: 700;
                    font-size: 21px;
                    margin: 0px;
                    color: color-mix(in srgb, var(--color-purple) 30%, black 70%)
                }

                .money {
                    align-self: flex-start;
                    font-family: 'Lexend';
                    font-weight: 600;
                    font-size: 23px;
                    margin: 0px;
                    text-wrap: nowrap;
                    background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                    border: 2px color-mix(in srgb, var(--color-purple) 20%, transparent 80%) solid;
                    padding: 2px 7px;
                    border-radius: 3px;
                }
            }

            .prob-selector {
                display: flex;
                flex-direction: column;
                gap: 10px;
                overflow-y: auto;
                min-height: 0px;
                background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                padding: 20px;
                box-sizing: border-box;
                width: 100%;
                flex-grow: 1;
                border-right: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                // border-right: 1px var(--color-purple) solid;

                .separator {
                    width: 100%;
                    border-top: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                    margin: 20px 0px;
                }

                &::-webkit-scrollbar {
                    display: none;
                }
                
            }
            
            .controls {
                background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
                border-top: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                border-right: 1px color-mix(in srgb, var(--color-purple) 30%, transparent 70%) solid;
                padding: 20px 20px;
                display: flex;
                justify-content: space-between;
                align-items: start;

                p {
                    font-family: 'Lexend';
                    font-size: 16px;
                    font-weight: 600;
                    margin: 0px;
                }

                i {
                    color: black;
                }

                .buy {
                    box-sizing: border-box;
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    align-items: flex-start;
                    gap: 10px;

                    .anchor {
                        position: relative;

                        .tooltip {
                            display: none;
                            position: absolute;
                            top: 50%;
                            left: 30px;
                            transform: translateY(-50%);
                            font-size: 12px;
                            font-family: 'Lexend';
                            color: color-mix(in srgb, var(--color-purple) 40%, black 60%);
                            padding: 5px 10px;
                            border-radius: 3px;
                            box-shadow: 0px 0px 5px 0px #00000030;
                            background-color: color-mix(in srgb, var(--color-purple) 15%, white 85%);
                            text-wrap: nowrap;

                            &.visible {
                                display: block;
                            }
                        }
                    }
                }

                .sell {
                    .anchor {
                        position: relative;

                        .tooltip {
                            display: none;
                            position: absolute;
                            top: 50%;
                            left: 30px;
                            transform: translateY(-50%);
                            font-size: 15px;
                            font-weight: 400;
                            font-family: 'Lexend';
                            color: color-mix(in srgb, var(--color-purple) 40%, black 60%);
                            padding: 5px 10px;
                            border-radius: 3px;
                            box-shadow: 0px 0px 5px 0px #00000030;
                            background-color: color-mix(in srgb, var(--color-purple) 15%, white 85%);
                            text-wrap: nowrap;

                            .bold {
                                font-weight: 700;
                            }

                            &.visible {
                                display: block;
                            }
                        }
                    }
                }
            }
        }
        
        .content {
            display: flex;
            flex-grow: 1;
            .prob {
                flex-grow: 1;
                display: flex;
                flex-direction: column;
                padding: 20px;
                box-sizing: border-box;

            }

            .chat {
                display: flex;
                flex-direction: column;
            }
        }
    }


</style>