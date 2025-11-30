<script lang="ts">
    import { onMount } from "svelte";
    import skrat from "$lib/assets/images/SKRAT.svg";

    const DT = 1/60;

    const collors: string[] = ["#9500EB", "#EB0072", "#EB6E00", "#EBAD00"];

    let width = $state(0);
    let height = $state(0);

    let canvas: HTMLCanvasElement;
    let ctx: CanvasRenderingContext2D | null;

    let xOftset = 0;
    let xOftsetTarget = 0;

    let creditsHovered = $state(false);

    onMount(() => {
        ctx = canvas.getContext("2d");

        animationFrame();
    });

    type Vector2 = {
        x: number,
        y: number
    }

    type FallingCube = {
        position: Vector2,
        speed: Vector2,
        orientation: number,
        angularSpeed: number,
        color: number,
        size: number,
        paralaxCof: number
    }

    let fallingCubes: FallingCube[] = [];

    function addRandomCube() {
        let position: Vector2 = { x: Math.random() * width, y: -10 };
        let speed: Vector2 = { x: (Math.random()*2 - 1) * 50, y: 0 }
        let orientation: number = Math.random()*360;
        let angularSpeed: number = (Math.random()*2 - 1) * 90;
        let color: number = Math.floor(Math.random()*4);
        let size: number = Math.random() * 15 + 5;
        let paralaxCof: number = 0.2 + Math.random() * 0.8;
        fallingCubes.push({ position, speed, orientation, angularSpeed, color, size, paralaxCof });
    }

    function updateCubes() {
        for (let cube of fallingCubes) {
            cube.speed.y += 1000 * DT * DT;

            cube.position.x += cube.speed.x * DT;
            cube.position.y += cube.speed.y * DT;
            cube.orientation += cube.angularSpeed * DT;
        }

        fallingCubes = fallingCubes.filter(cube => cube.position.y < height + 10);
    }

    function renderCubes() {
        if(!ctx) return;

        ctx.clearRect(0, 0, width, height);
        
        for(let cube of fallingCubes) {
            ctx.fillStyle = collors[cube.color];
            ctx.save();
            ctx.translate(cube.position.x + xOftset * cube.paralaxCof, cube.position.y);
            ctx.rotate(cube.orientation * Math.PI / 180);

            ctx.beginPath();
            ctx.roundRect(-cube.size/2, -cube.size/2, cube.size, cube.size, 3);
            ctx.fill();
            
            ctx.restore();

        }
    }

    function updateParalax() {
        xOftset += (xOftsetTarget - xOftset) * 0.05;
    }



    function animationFrame() {
        requestAnimationFrame(animationFrame);

        if (Math.random() < 0.2) {
            addRandomCube();
        }

        updateParalax();
        updateCubes();
        renderCubes();
    }


</script>

<main bind:clientWidth={width} bind:clientHeight={height} onmousemove={(e) => {
    xOftsetTarget = (e.clientX - width) / 5;
}}>
    <div class="title">
        <h1>Vyčkejte na začátek soutěže</h1>
    </div>
    <div class="credits">
        <img src={skrat} alt="credits to skrat team" onmouseenter={() => { creditsHovered = true; }} onmouseleave={() => { creditsHovered = false; }}>
        <div class="content" class:hovered={creditsHovered}>
            <h3 class="title">Credits:</h3>
            <p style="transition-delay: 0.05s"><span class="bold">Ondřej Urban:</span> frontend developer, designer</p>
            <p style="transition-delay: 0.10s"><span class="bold">Rostislav Kozlík:</span> backend developer</p>
            <p style="transition-delay: 0.15s"><span class="bold">Šimon Vecka:</span> on-site round hardware engineer, photographer</p>
            <p style="transition-delay: 0.20s"><span class="bold">Tomáš Jodl:</span> on-site round hardware engineer, photographer</p>
        </div>
    </div>
    <canvas bind:this={canvas} width={width} height={height}></canvas>
</main>

<style lang="scss">
    main {
        position: relative;
        display: flex;
        flex-direction: column;
        align-items: center;
        height: 100%;
        background-color: white;
        
        .title {
            position: absolute;
            top: 50px;

            h1 {
                font-family: 'Fredoka';
                margin: 0px;

                font-size: 50px;
                font-weight: 600;
                backdrop-filter: blur(3px);
                background-color: #F0F0F030;
                box-sizing: border-box;
                box-shadow: 0px 0px 10px #F0F0F0;
                border-radius: 5px;
                padding: 5px 10px;
                color: color-mix(in srgb, var(--color-purple) 40%, black 60%);
            }

        }

        .credits {
            user-select: none;
            position: absolute;
            bottom: 30px;
            width: 100%;
            display: flex;
            justify-content: center;

            img {
                width: 250px;
                opacity: 0.03;
                cursor: pointer;
            }

            .content {
                // width: 100%;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                position: absolute;
                top: -20px;
                transform: translateY(-100%);
                gap: 10px;
                padding: 10px 20px;

                border-radius: 5px;
                transition: all 0.5s cubic-bezier(0.215, 0.610, 0.355, 1);

                h3.title {
                    all: unset;
                    font-family: 'Fredoka';
                    font-size: 25px;
                    font-weight: 600;
                    color: color-mix(in srgb, var(--color-purple) 30%, #444444 70%);
                    transform: TranslateY(30px);
                    opacity: 0;
                    transition: all 0.5s cubic-bezier(0.55, 0.055, 0.675, 0.19);
                    // transition: all 0.5s cubic-bezier(0.215, 0.610, 0.355, 1);
                }

                p {
                    font-family: 'Fredoka';
                    font-size: 16px;
                    font-weight: 500;
                    color: color-mix(in srgb, var(--color-purple) 10%, gray 90%);
                    margin: 0px;
                    transform: TranslateY(30px);
                    opacity: 0;
                    transition: all 0.5s cubic-bezier(0.55, 0.055, 0.675, 0.19);

                    .bold {
                        font-weight: 600;
                        color: color-mix(in srgb, var(--color-purple) 50%, black 50%);
                    }
                }

                &.hovered {
                    backdrop-filter: blur(3px);
                    background-color: #F0F0F030;
                    box-sizing: border-box;
                    box-shadow: 0px 0px 10px #F0F0F0;

                    p, h3.title {
                        transition: all 0.5s cubic-bezier(0.215, 0.610, 0.355, 1);

                        opacity: 1;
                        transform: translateY(0px);
                    }
                }
            }
        }
    }
</style>