<script lang="ts">
    import Button from "$lib/components/general/Button.svelte";
    
    import QRreader from "$lib/components/admin/QRreader.svelte";
    // import LoadingAnimationColor from "shared";
    import { LoadingAnimation }  from "shared";
    // import { LoadingAminationColor } from "$lib/components/general/LoadingAnimationColor.svelte"

    let problemScanned: boolean = $state(false);
    let problemReady: boolean = $state(false);

    let answerVisible: boolean = $state(false);

    let confirmationType: "sell" | "accept" | undefined = $state(undefined);

    let scannedProb: {
        id: string,
        money: number,
        answer: string,
        state: "bought" | "sold" | "solved" | "free"
    } | undefined = $state({id: "idididididididididi", money: 567, answer: "Tohle je odpovwed na priklad", state: "solved"});

    async function handleSubmit(type: "sell" | "accept") {
        if (type == "accept") {
            problemReady = false;
            try {
                let response = await fetch(`https://strela-vlna.gchd.cz/api/gradeprob?id=${scannedProb?.id}`);

            } catch (error: any) {
                alert(error);
            } finally {
                problemScanned = false;
            }
        } else if (type == "sell") {
            problemReady = false;
            try {
                let response = await fetch(`https://strela-vlna.gchd.cz/api/sellprob?id=${scannedProb?.id}`);

            } catch (error: any) {
                alert(error);
            } finally {
                problemScanned = false;
            }
        }
    }

    async function handleScan(text: string) {
        if (problemScanned) return;
        problemScanned = true;
        
        try {
            let response: any = await (await fetch(`https://strela-vlna.gchd.cz/api/getprob?id=${text}`)).json();
            console.log(response);

            
            // if (response.status != 200){
            //     console.log("sdjhfkahj");
            //     throw response.message;
            // }
            
            scannedProb = {
                id: text,
                money: response.money,
                answer: response.answer,
                state: response.state
            }
            
            problemReady = true;
        
        } catch (error: any) {
            alert(error);
            problemScanned = false;

        }
    }

</script>

<main>
    <QRreader onScan={handleScan}/>
    
    <div class="controls">
        {#if confirmationType !== undefined}
        <div class="confirmation-wrapper">
            <div class="content">
                {#if confirmationType == "sell"}
                    <h2>Do you really want to <span class="bold">sell</span> this problem?</h2>
                {:else if confirmationType == "accept"}
                    <h2>Accept the solution?</h2>
                {/if}
    
                <div class="controls">
                    <Button theme="pink" disabled={false} onclick={() => {
                        confirmationType = undefined;
                    }}>
                        <i class="fa-solid fa-xmark"></i>
                        <p>Cancel</p>
                    </Button>
                    <Button theme="yellow" disabled={false} onclick={() => {
                        let cType = confirmationType;
                        confirmationType = undefined;
                        handleSubmit(cType as "sell" | "accept");
                    }}>
                        <i class="fa-solid fa-check"></i>
                        <p>Submit</p>
                    </Button>
                </div>
            </div>
        </div>
        {/if}
        <div class="content">
            <div class="top-row">
                <Button theme="pink" disabled={!problemReady} onclick={() => {
                    if (!problemReady) return;
                    confirmationType = "sell";
                }}>
                    <i class="fa-solid fa-trash-can"></i>
                    <p>Sell</p>
                </Button>
                {#if problemReady}
                    <div class="money-state-wrapper">
                        <h2 class="money">{scannedProb?.money ?? "-"} DC</h2>
                        <div class="state"
                            class:bought={scannedProb?.state=="bought"}
                            class:sold={scannedProb?.state=="sold"}
                            class:solved={scannedProb?.state=="solved"}
                            class:free={scannedProb?.state=="free"}></div>
                    </div>
                    <div class="back">
                        <Button theme="purple" disabled={!problemReady} onclick={() => {
                            if (!problemReady) return;
    
                            problemReady = false;
                            problemScanned = false;
                        }}>
                            <i class="fa-solid fa-arrow-turn-up"></i>
                        </Button>
                    </div>
                {/if}
                {#if !problemReady && problemScanned}
                    <div class="animation-wrapper">
                        <LoadingAnimation color="black"/>
                    </div>
                {/if}


            </div>
            <div class="bottom-row">
                <Button theme="yellow" disabled={!problemReady} onclick={() => {
                    if (!problemReady) return;
                    confirmationType = "accept";
                }}>
                    <i class="fa-regular fa-square-check"></i>
                    <p>Accept</p>
                </Button>
                <div class="answer">
                    <Button theme="orange" disabled={!problemReady}
                        ontouchstart={() => {
                            if (!problemReady) return;
                            answerVisible = true;
                        }}
                        ontouchcancel={() => {
                            if (!problemReady) return;
                            answerVisible = false;
                        }}
                        ontouchend={() => {
                            if (!problemReady) return;
                            answerVisible = false;
                        }}

                        onmousedown={() => {
                            if (!problemReady) return;
                            answerVisible = true;
                        }}
                        onmouseup={() => {
                            if (!problemReady) return;
                            answerVisible = false;
                        }}
                        >
                        <i class="fa-regular fa-circle-question"></i>
                        <p>Answer</p>
                    </Button>
                    {#if answerVisible}
                        <div class="answer-text">
                            <p>{scannedProb?.answer ?? ""}</p>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
</main>

<style lang="scss">
    main {
        height: 100vh;
        width: 100vw;

        .controls {
            position: absolute;
            top: 0px;
            left: 0px;
            width: 100%;
            height: 100%;

            .confirmation-wrapper {
                position: absolute;
                padding-inline: 30px;
                z-index: 100;
                background-color: #00000040;
                height: 100%;
                width: 100%;
                box-sizing: border-box;
                display: flex;
                flex-direction: column;
                justify-content: center;

                .content {
                    background-color: white;
                    height: fit-content;
                    border-radius: 5px;
                    padding: 10px;

                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    align-items: center;

                    h2 {
                        font-family: 'Lexend';
                        font-weight: 500;
                        width: 100%;
                        text-align: center;
                        font-size: 25px;
                    }

                    .bold {
                        font-weight: 800;
                        font-size: 28px;
                    }

                    .controls {
                        position: relative;
                        display: flex;
                        flex-direction: row;
                        align-items: center;
                        justify-content: space-between;
                        padding-top: 30px;

                        p {
                            font-family: 'Lexend';
                            font-size: 18px;
                            font-weight: 600;
                            margin: 0px;
                        }

                        i {
                            font-size: 22px;
                        }
                    }

                }
            }

            .content {
                display: flex;
                flex-direction: column;
                align-items: flex-start;
                justify-content: space-between;
                height: 100%;
                padding: 0px 15px;
                
                .top-row {
                    background-color: white;
                    box-shadow: 0px 0px 10px #2c2c2c52;
                    padding: 10px 10px;
                    box-sizing: border-box;
                    width: 100%;
                    border-bottom-left-radius: 5px;
                    border-bottom-right-radius: 5px;

                    display: flex;
                    justify-content: space-between;
                    align-items: center;

                    i {
                        font-size: 18px;
                    }

                    p {
                        font-family: 'Lexend';
                        margin: 0px;
                        font-size: 20px;
                    }

                    .animation-wrapper {
                        width: 35px;
                        height: 35px;
                        
                    }

                    .money-state-wrapper {
                        display: flex;
                        width: fit-content;
                        flex-direction: column;
                        gap: 0px;
                        
                        .state {
                            width: 100%;
                            border: none;
                            border-bottom: 5px solid;

                            &.free {
                                border-color: var(--color-purple);
                            }
                            &.bought {
                                border-color: var(--color-yellow);
                            }
                            &.solved {
                                border-color: var(--color-orange);
                                border-bottom-style: dashed;
                            }
                            &.sold {
                                border-color: var(--color-pink);
                            }

                        }

                        h2.money {
                            font-family: 'Lexend';
                            margin: 0px;
                            font-weight: 700;
                            font-size: 30px;
                            padding-inline: 5px;
                        }
                    }

                    .back {
                        i {
                            font-size: 25px;
                            transform: rotate(-90deg);
                        }
                    }
                }

                .bottom-row {
                    background-color: white;
                    box-shadow: 0px 0px 10px #2c2c2c52;
                    padding: 10px 10px;
                    box-sizing: border-box;
                    width: 100%;
                    border-top-left-radius: 5px;
                    border-top-right-radius: 5px;

                    display: flex;
                    align-items: center;
                    justify-content: space-between;
                    flex-direction: row;

                    i {
                        font-size: 20px;
                    }
                
                    p {
                        font-family: 'Lexend';
                        margin: 0px;
                        font-size: 18px;
                        font-weight: 600;
                    }

                    .answer {
                        position: relative;

                        p {
                            font-family: 'Lexend';
                            margin: 0px;
                            font-size: 18px;
                            font-weight: 600;
                            user-select: none;
                        }

                        .answer-text {
                            position: absolute;
                            display: flex;
                            padding: 5px 10px;
                            background-color: white;
                            border-radius: 5px;
                            top: 0px;
                            left: 50%;
                            box-sizing: border-box;

                            transform: translate(-50%, calc(-100% - 30px));

                            p {
                                font-size: 16px;
                                font-weight: 400 !important;
                                margin: 0px;
                            }
                        }
                    }
                }
            }

        }
    }
</style>