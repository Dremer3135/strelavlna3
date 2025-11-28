<script lang="ts">
    import { onMount } from "svelte";
    type Segment = "A" | "B" | "C" | "D" | "E" | "F" | "G";

    let { digit, color }: { digit: number, color: string } = $props();


    let segmentCombinations: Segment[][] = [
        [
            "A", "B", "C", "D", "E", "F"
        ],[
            "B", "C"
        ],[
            "A", "B", "G", "E", "D"
        ],[
            "A", "B", "G", "C", "D"
        ],[
            "A", "F", "G", "B", "C"
        ],[
            "F", "G", "C", "D", "A"
        ],[
            "F", "E", "D", "C", "G"
        ],[
            "A", "B", "C"
        ],[
            "A", "F", "G", "C", "D", "E", "B"
        ],[
            "A", "F", "G", "B", "C"
        ]

    ];

    $effect(() => {
        if(!ctx) {
            drawDigit(digit);
        }
    })




    let canvas: HTMLCanvasElement;
    let ctx: CanvasRenderingContext2D | null;
    
    let height = $state(0);
    let width = $derived(height / 9 * 5);
    let pointSize = $derived(height / 9);

    onMount(() => {
        ctx = canvas.getContext('2d');
    });

    // $effect(() => {
    //     if (ctx) {
    //         ctx.clearRect(0, 0, width, height);


    //         ctx.fillStyle=color;
    //         ctx.beginPath();
    //         ctx.roundRect(0, 0, 10, 50, 3);
    //         ctx.fill();
    //     }
    // });

    function drawDot(x: number, y: number) {
        if (!ctx) return;

        ctx.fillStyle=color;
        ctx.beginPath();
        ctx.roundRect(x*pointSize, y*pointSize, pointSize, pointSize, pointSize/5);
        ctx.fill();
    }

    function drawSegment(segment: Segment) {
        if (segment === "A") {
            drawDot(1, 0);
            drawDot(2, 0);
            drawDot(3, 0);
        }
        else if (segment === "B") {
            drawDot(4, 1);
            drawDot(4, 2);
            drawDot(4, 3);
        }
        else if (segment === "C") {
            drawDot(4, 5);
            drawDot(4, 6);
            drawDot(4, 7);
        }
        else if (segment === "D") {
            drawDot(1, 8);
            drawDot(2, 8);
            drawDot(3, 8);
        }
        else if (segment === "E") {
            drawDot(0, 5);
            drawDot(0, 6);
            drawDot(0, 7);
        }
        else if (segment === "F") {
            drawDot(0, 1);
            drawDot(0, 2);
            drawDot(0, 3);
        }
        else if (segment === "G") {
            drawDot(1, 4);
            drawDot(2, 4);
            drawDot(3, 4);
        }
    }

    function drawDigit(digit: number) {
        // if (!ctx) return;
        
        // ctx.clearRect(0, 0, width, height);

        for(let i = 0; i < segmentCombinations[digit].length; i++) {
            let segment = segmentCombinations[digit][i];

            setTimeout(() => {
                drawSegment(segment);
            }, i*40);
        }
    }

</script>

<main bind:clientHeight={height}>
    <canvas bind:this={canvas} width={width} height={height}></canvas>
</main>

<style lang="scss">
    main {
        height: 100%;
        width: fit-content;
    }
</style>