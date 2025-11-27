<script lang="ts">
    import type { Prob } from "$lib/types";
    import type { MessageType } from "$lib/types";
    import Message from "./Message.svelte";

    let { prob, send }: { prob: Prob, send: (message: MessageType) => void } = $props();

    let inputMessage: string = $state("");

    function handleSend() {
        if (inputMessage.length == 0) return;
        
        send({
            origin: "sent",
            type: "message",
            value: inputMessage,
        });
        
        
        inputMessage = "";
    }

</script>

<main>
    <div class="title">
        <h2>Chat</h2>
    </div>
    <div class="content">
        {#each prob.chat as message}
            <Message message={message} />
        {/each}    
    </div>
    <div class="controls">
        <input type="text" placeholder="Zpráva" bind:value={inputMessage} onkeydown={(e) => {
            if (e.key =="Enter") handleSend();
        }}>
        <button class="submit" class:ready={ inputMessage.length > 0 } onclick={() => { handleSend(); }}>
            <i class="fa-solid fa-paper-plane"></i>
        </button>   
    </div>
</main>

<style lang="scss">
    main {
        background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
        border-left: 1px color-mix(in srgb, var(--color-pink) 30%, transparent 70%) solid;
        width: 350px;
        padding: 20px;
        box-sizing: border-box;
        flex-grow: 1;
        display: flex;
        flex-direction: column;

        .title {
            h2 {
                font-family: 'Lexend';
                font-size: 30px;
                font-weight: 600;
                margin: 0px;
                text-align: center;
                padding-bottom: 10px;
                border-bottom: 1px color-mix(in srgb, var(--color-pink) 30%, transparent 70%) solid;
                color: color-mix(in srgb, var(--color-pink) 30%, black 70%);
            }
        }

        .content {
            min-height: 0px;
            overflow-y: auto;
            padding: 20px 0px;
            display: flex;
            flex-direction: column;
            gap: 10px;
            flex-grow: 1;
        }

        .controls {
            padding-bottom: 20px;
            box-sizing: border-box;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 10px;

            input {
                all: unset;
                cursor: text;
                font-family: 'Lexend';
                font-size: 18px;
                font-weight: 600;
                padding: 7px 10px;
                border: 3px color-mix(in srgb, var(--color-pink) 30%, white 70%) solid;
                border-radius: 3px;
                background-color: color-mix(in srgb, #FAFAFA 50%, transparent 50%);

                &:focus {
                    border-color: var(--color-pink);
                    background-color: #FAFAFA;
                }
            }

            .submit {
                all: unset;
                cursor: pointer;
                height: 100%;
                min-width: 43px;
                border: 3px color-mix(in srgb, var(--color-pink) 30%, white 70%) solid;
                box-sizing: border-box;
                display: flex;
                justify-content: center;
                align-items: center;
                border-radius: 5px;

                &:hover {
                    background-color: color-mix(in srgb, var(--color-pink) 30%, white 70%);
                }

                &.ready {
                    border-color: var(--color-pink);

                    &:hover {
                        background-color: var(--color-pink);
                        color: white;
                    }
                }

            }
        }
    }





</style>