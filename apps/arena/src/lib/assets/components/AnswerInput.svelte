<script lang="ts">
    import { isProbSolved } from "$lib/utils";

    let { prob, submitAnswer, disabled, placeholder, onPaste }: { prob: any, submitAnswer: (answer: string) => void, disabled: boolean, placeholder: string, onPaste?: (data: string) => void } = $props();

    let answer: string = $state("");

    $effect(() => {
        prob;
        answer = "";
    });

</script>

<main>
    <input type="text" placeholder={placeholder} class:disabled={disabled} bind:value={answer} class:placeholder-invisible={isProbSolved(prob)} onkeydown={(e) => {
        if (e.key == "Enter") {
            if (disabled) return;
            submitAnswer(answer);
            answer = "";
        }
    }}
    onpaste={(e) => {
        const pastedText = e.clipboardData?.getData('text');
        if (pastedText && onPaste) {
            onPaste(pastedText);
        }
    }}
    disabled={disabled}/>
</main>

<style lang="scss">
    main {
        width: 100%;
        display: flex;
        padding: 20px;
        box-sizing: border-box;

        input {
            all: unset;
            cursor: text;
            flex-grow: 1;
            box-sizing: border-box;
            background-color: #F0F0F0;
            border: 3px #CCCCCC solid;
            padding: 7px 15px;
            font-family: 'Lexend';
            font-weight: 600;
            font-size: 18px;
            border-radius: 3px;

            &:focus:not(.disabled) {
                border-color: var(--color-purple);
                background-color: color-mix(in srgb, var(--color-purple) 5%, transparent 95%);
            }

            &.disabled {
                background-color: #FAFAFA;
                border-color: #EEEEEE;
                cursor: not-allowed;
            }

            &.placeholder-invisible::placeholder {
                opacity: 0;
            }


        }
    }
</style>