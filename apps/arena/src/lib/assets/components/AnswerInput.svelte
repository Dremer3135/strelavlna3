<script lang="ts">
    let { submitAnswer, disabled, placeholder, onPaste }: { submitAnswer: (answer: string) => void, disabled: boolean, placeholder: string, onPaste?: (data: string) => void } = $props();

    let answer: string = $state("");

</script>

<main>
    <input type="text" placeholder={placeholder} class:disabled={disabled} bind:value={answer} onkeydown={(e) => {
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
    }}>
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
            background-color: #FAFAFA;
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


        }
    }
</style>