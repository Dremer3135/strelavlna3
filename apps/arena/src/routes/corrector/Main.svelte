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
	import { getProbLastAnswer, hasProbChat, isProbSolved, isProbUngraded } from '$lib/utils';
    import ProbSelectCorrector from '$lib/assets/components/ProbSelectCorrector.svelte';
    import Grader from '$lib/assets/components/Grader.svelte';

	let {
		focus,
		chat,
    start,
	}: {
		focus: (id: string) => void;
		chat: (probId: string, message: Omit<MessageType, 'sentTime'>) => void;
        start: () => void;
	} = $props();

	const difficulties = ['A', 'B', 'C'];

</script>

<main>
  <button onclick={start}>Start</button>
	<div class="left-panel">
		<!-- <div class="team-stats">
			<h2 class="name">{$currentState.teamName}</h2>
			<h3 class="money">{$currentState.money} DC</h3>
		</div> -->
		<div class="prob-selector">
			{#each Object.values($probs).filter((prob) => isProbUngraded(prob)) as prob}  <!--  Ungraded  -->
				<ProbSelectCorrector
					prob={prob}
					onSelect={() => {
						focus(prob.id);
					}}
				/>
			{/each}
			<div class="separator"></div>
			{#each Object.values($probs).filter((prob) => hasProbChat(prob) && !isProbUngraded(prob)) as prob}  <!--  Chat started  -->
				<ProbSelectCorrector
					prob={prob}
					onSelect={() => {
						focus(prob.id);
					}}
				/>
			{/each}
			<div class="separator"></div>
			{#each Object.values($probs).filter((prob) => !isProbSolved(prob) && !hasProbChat(prob) && !isProbUngraded(prob)) as prob}  <!--  Newly buyed  -->
				<ProbSelectCorrector
					prob={prob}
					onSelect={() => {
						focus(prob.id);
					}}
				/>
			{/each}
			<div class="separator"></div>
			{#each Object.values($probs).filter((prob) => isProbSolved(prob)) as prob}  <!--  Solved  -->
				<ProbSelectCorrector
					prob={prob}
					onSelect={() => {
						focus(prob.id);
					}}
				/>
			{/each}
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
				<Grader answer={getProbLastAnswer($focusedProb) ?? ""} prob={$focusedProb} onAccept={() => {
						chat($focusedProb.id, {
							origin: "sent",
							type: "grade",
							value: "correct"
						});
					}} 
					onReject={() => {
						chat($focusedProb.id, {
							origin: "sent",
							type: "grade",
							value: "incorrect"
						});
					}}
				/>
			</div>
			<div class="chat">
				<Chat
					type="corrector"
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
			min-height: 0px;

            .prob {
                flex-grow: 1;
                display: flex;
                flex-direction: column;
                padding: 20px;
                box-sizing: border-box;

            }

            .chat {
				min-height: 0px;
                display: flex;
                flex-direction: column;
            }
        }
    }


</style>
