<script lang="ts">
    import type { MessageType } from "$lib/types";
    let { message }: { message: MessageType } = $props();
</script>

<main class:recieved={message.origin === "sent"}>
    {#if message.type == "message"}
    <div class="message">
        <div class="content">
            <p>{message.value}</p>
        </div>
    </div>
    {:else if message.type == "grade"}
    <div class="grade">
        <div class="content">
            <p>{message.value}</p>
        </div>
    </div>
    {/if}
</main>

<style lang="scss">
    main {
        border: 3px var(--color-pink) solid;
        border-radius: 5px;
        background-color: color-mix(in srgb, var(--color-pink) 5%, transparent 95%);
        padding: 10px;
        box-sizing: border-box;
        margin-right: 30px;
        
        .content {
            p {
                font-family: 'Lexend';
                font-size: 16px;
                font-weight: 500;
                margin: 0px;
            }
        }
        
        &.recieved {
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



</style>