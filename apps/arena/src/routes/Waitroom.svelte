<script lang="ts">
    import { onMount } from "svelte";
    import skrat from "$lib/assets/images/SKRAT.svg";

    const DT = 1/60;

    const collors: string[] = ["#9500EB", "#EB0072", "#EB6E00", "#EBAD00"];

    let width: number;
    let height: number;

    let canvas: HTMLCanvasElement;
    let ctx: CanvasRenderingContext2D | null;

    let xOftset = 0;
    let xOftsetTarget = 0;

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
        let paralaxCof: number = Math.random();
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
        <h1>Vyčkejte než začne soutěž</h1>
    </div>
    <div class="credits">
        <img src={skrat} alt="credits to skrat team">
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
                color: color-mix(in srgb, var(--color-purple) 40%, black 60%)
            }

        }

        .credits {
            position: absolute;
            bottom: 30px;

            img {
                width: 250px;
                opacity: 0.03;
            }
        }
    }
</style>