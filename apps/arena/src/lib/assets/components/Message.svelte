<script lang="ts">
    import type { MessageType } from "$lib/types";
    let { message, type }: { message: MessageType, type: "player" | "corrector" } = $props();
    
    let sentTimeString = $derived(
        `${String(message.sentTime.getHours()).padStart(2, '0')}:` +
        `${String(message.sentTime.getMinutes()).padStart(2, '0')}`);
</script>

<main>
    {#if message.type == "message"}
    <div class="message" class:sent={message.origin === "sent"}>
        <p class="time">{sentTimeString}</p>
        <div class="content">
            <p>{message.value}</p>
        </div>
    </div>
    {:else if message.type == "grade"}
    <div class="grade" class:sent={message.origin === "sent"}>
        <p class="time">{sentTimeString}</p>
        <div class="content">
            {#if message.value == "correct"}
            <p>Odpověď byla uznána!</p>
            <i class="fa-solid fa-square-check"></i>
            {:else}
            <p>Odpověď byla zamítnuta!</p>
            <i class="fa-solid fa-square-xmark"></i>
            {/if}
        </div>
    </div>
    {:else if message.type == "answer"}
    <div class="answer" class:sent={message.origin === "sent"}>
        <p class="time">{sentTimeString}</p>
        <div class="content">
            <p>Odpověděli jste: <span class="bold">{message.value}</span></p>
        </div>
    </div>
    {:else if message.type == "copy" && type == "corrector"}
    <div class="copy-paste" class:sent={message.origin === "sent"}>
        <p class="time">{sentTimeString}</p>
        <div class="content">
            <p>Hráč zkopíroval: <span class="bold">{message.value}</span></p>
        </div>
    </div>
    {:else if message.type == "paste" && type == "corrector"}
    <div class="copy-paste" class:sent={message.origin === "sent"}>
        <p class="time">{sentTimeString}</p>
        <div class="content">
            <p>Hráč vložil: <span class="bold">{message.value}</span></p>
        </div>
    </div>
    {:else if message.type == "window-focus" && type == "corrector"}
    <div class="window-focus" class:sent={message.origin === "sent"}>
        <p class="time">{sentTimeString}</p>
        <div class="content">
            {#if message.value.split(" ")[0] == "focused"}
                <p>{message.value.split(" ")[1]} přišel</p>
                <i class="fa-solid fa-right-to-bracket"></i>
            {:else}
                <p>{message.value.split(" ")[1]} odešel</p>
                <i class="fa-solid fa-right-from-bracket"></i>
            {/if}
        </div>
    </div>
    {/if}
</main>

<style lang="scss">
    main {
        .time {
            position: absolute; 
            font-family: 'Lexend';
            font-size: 14px;
            font-weight: 500;
            color: gray;
            right: 0px;
            transform: translateX(calc(100% + 10px));
            top: 0px;
            margin: 0px;
            margin-top: 2px;
        }
        
        .sent {
            margin-right: 0px !important;
            margin-left: 50px !important;
    
            .time {
                right: unset;
                left: 0px;
                transform: translateX(calc(-100% - 10px));
            }
        }

        .message {
            position: relative;
            border: 3px var(--color-pink) solid;
            border-radius: 5px;
            background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
            padding: 10px;
            box-sizing: border-box;
            margin-right: 50px;


            .content {
                p {
                    font-family: 'Lexend';
                    font-size: 16px;
                    font-weight: 500;
                    margin: 0px;
                }
            }

            &.sent {
                margin-right: 0px;
                margin-left: 30px;
                background-color: var(--color-pink);
            
                .content {
                    p {
                        color: white;
                        font-weight: 600;
                    }
                }

            }
        }

        .grade {
            position: relative;
            border: 3px var(--color-purple) solid;
            border-radius: 5px;
            padding: 10px;
            box-sizing: border-box;
            margin-right: 50px;

            background-image: repeating-linear-gradient(
                -45deg,
                color-mix(in srgb, var(--color-purple) 3%, white 97%),
                color-mix(in srgb, var(--color-purple) 0%, white 100%) 10px,
                color-mix(in srgb, var(--color-purple) 10%, white 90%) 10px,
                color-mix(in srgb, var(--color-purple) 6%, white 94%) 20px,
            );
            
            
            .content {
                display: flex;
                gap: 10px;
                align-items: center;
                justify-content: space-between;
                
                p {
                    font-family: 'Lexend';
                    font-size: 16px;
                    font-weight: 700;
                    margin: 0px;
                    color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
                    
                }

                i {
                    font-size: 22px;
                    color: color-mix(in srgb, var(--color-purple) 65%, black 35%)
                }
            }

            &.sent {
            }
        }

        .answer {
            position: relative;
            border: 3px var(--color-pink) solid;
            border-radius: 5px;
            padding: 10px;
            box-sizing: border-box;
            margin-right: 50px;

            background-image: repeating-linear-gradient(
                -45deg,
                color-mix(in srgb, var(--color-pink) 3%, white 97%),
                color-mix(in srgb, var(--color-pink) 0%, white 100%) 10px,
                color-mix(in srgb, var(--color-pink) 10%, white 90%) 10px,
                color-mix(in srgb, var(--color-pink) 6%, white 94%) 20px,
            );
            
            
            .content {
                p {
                    font-family: 'Lexend';
                    font-size: 16px;
                    font-weight: 400;
                    margin: 0px;

                    .bold {
                        font-weight: 700;
                    }
                }
            }
            
            // &.sent {
            //     background-image: none;
            //     background-color: color-mix(in srgb, var(--color-pink) 15%, white 85%);
            // }
        }
    }



</style>