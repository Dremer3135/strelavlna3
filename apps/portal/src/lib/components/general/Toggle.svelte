<script lang="ts">
    type OnChangeFunction = (value: boolean) => void;

    let { value, onChangeFunction }: { value: boolean, onChangeFunction: OnChangeFunction } = $props();



    let sliderWidth: number = $state(0);
    let sliderHeight: number = $state(0);

    $effect(() => {
        console.log(sliderWidth);
    });


</script>

<main>
    <button style="
    width: {sliderWidth*2}px;
    height: {sliderHeight}px;
    "
    aria-label="toggle"
    class:active={ value }
    onclick={() => { onChangeFunction(!value); }}
    >
        <div class="slider-wrapper" bind:clientHeight={ sliderHeight } bind:clientWidth={ sliderWidth }>
            <div class="slider"></div>
        </div>
    </button>
</main>

<style lang="scss">
    main {

        button {
            all: unset;
            position: relative;
            border-radius: 6px;
            border: 3px lightgray solid;
            cursor: pointer;
            
            background-color: #FAFAFA;
            transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
            
            .slider-wrapper {
                position: absolute;
                top: 0px;
                left: 0px;
                width: 20px;
                height: 20px;
                padding: 3px;
                transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
                
                .slider{
                    width: 100%;
                    height: 100%;
                    box-sizing: border-box;
                    background-color: lightgray;
                    border-radius: 3px;
                    box-shadow: 0px 0px 7px #00000030;
                    transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
                }
            }
            
            &.active {
                border-color: color-mix(in srgb, var(--color-pink) 90%, black 10%);
                background-color: color-mix(in srgb, var(--color-pink) 40%, transparent 60%);
                
                .slider-wrapper {
                    left: 50%;
                    
                    .slider {
                        background-color: color-mix(in srgb, var(--color-pink) 2%, white 98%);

                    }
                }

            }
        
        }



    }
</style>