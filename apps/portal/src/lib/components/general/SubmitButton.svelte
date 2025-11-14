<script lang="ts">
    let { isLoading = false, children, ...rest }: { isLoading: boolean, children: any } = $props();

    // let { class: className = '', theme  }: { class: any, theme : "yellow" | "orange" | "purple" | "pink", children: any } = $props();


    let childrenWidth = $state(0);


</script>

<button style="--children-width: {childrenWidth}px" class:loading={isLoading}>
    <span bind:clientWidth={childrenWidth}>{@render children()}</span>
    <i class="fa-solid fa-angle-right arrow"></i>
    <i class="fa-solid fa-rotate loading"></i>
</button>

<style lang="scss">
    button {
        all: unset;
        cursor: pointer;
        position: relative;
        box-sizing: border-box;
        background-color: black;
        border-radius: 4px;
        padding: 10px 20px;
        width: 100%;
        max-width: 400px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;

        
    
        span {
            font-family: 'Lexend';
            font-size: 17px;
            color: white;
            transform: translateX(0);

            transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
        }

        .arrow {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, calc(-50% + 1px));
            font-size: 17px;
            color: white;
            opacity: 0;
            transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
        }

        .loading {
            font-size: 20px;
            position: absolute;
            top: 50%;
            left: 20px;
            transform: translateY(-50%);
            animation: loading-animation 1s linear infinite;
            opacity: 0;
        }

        &.loading {
            .loading {
                opacity: 1;
            }
        }

        &:hover {
            span {
                transform: translateX(calc(-1 * (var(--children-width) / 4 + 0px)));
            }

            .arrow {
                transform: translate(calc(-50% + var(--children-width) / 4 + 20px), calc(-50% + 1px));
                opacity: 1;
            }
        }
    }

    @keyframes loading-animation {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
</style>