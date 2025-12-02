<script lang="ts">
let canvas: HTMLCanvasElement | null = null;
let loading_ctx: CanvasRenderingContext2D | null = null;

class Block {
    x: number;
    y: number;
    target_x: number;
    target_y: number;
    destroy_x: number;
    destroy_y: number;

    constructor(x: number, y: number, target_x: number, target_y: number, destroy_x: number, destroy_y: number) {
        this.x = x;
        this.y = y;
        this.target_x = target_x;
        this.target_y = target_y;
        this.destroy_x = destroy_x;
        this.destroy_y = destroy_y;
    }
}

class Pos {
    x: number;
    y: number;
    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
    }
}

class Void {
    x: number;
    y: number;
    destination_x: number;
    destination_y: number;

    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
        this.destination_x = x;
        this.destination_y = y;
    }
}


function voidIndex(void_: Void, list: Void[]): number {
    return list.findIndex((e) => e.x === void_.x && e.y === void_.y);
}

function uppdate_blocks(blocks: Block[], k: number) {
    for (let i = 0; i < blocks.length; i++) {
        const b = blocks[i];
        b.x += (b.target_x - b.x) * k;
        b.y += (b.target_y - b.y) * k;

        if (Math.abs(b.destroy_x - b.x) < 0.05 && Math.abs(b.destroy_y - b.y) < 0.05) {
            blocks.splice(i, 1);
            i--;
        } else if (Math.abs(b.target_x - b.x) < 0.01 && Math.abs(b.target_y - b.y) < 0.01) {
            b.x = b.target_x;
            b.y = b.target_y;
        }
    }
}

function draw_blocks() {
    if (!loading_ctx || !canvas) return;

    loading_ctx.clearRect(0, 0, canvas.width, canvas.height);
    for (const block of canvas_accual) {
        loading_ctx.beginPath();
        loading_ctx.roundRect(30 * block.x, 30 * block.y, spacing, spacing, corner);
        loading_ctx.fill();
    }
}

function blockIndex(block: Block, list: Block[]): number {
    return list.findIndex((b) => b.target_x === block.target_x && b.target_y === block.target_y);
}

function count_blocks(trgt_canvas: number[][]): number {
    let i = 0;
    for (let y = 0; y < 10; y++)
        for (let x = 0; x < 22; x++)
            if (trgt_canvas[y][x] === 0) i++;
    return i;
}

function change_screen(starting_canvas: number[][], target_canvas: number[][]) {
    let delta_N = count_blocks(target_canvas) - count_blocks(starting_canvas);
    let modified_canvas = starting_canvas.map(row => [...row]);

    if (delta_N !== 0) {
        if (delta_N > 0) {
            for (let y = 0; y < 10 && delta_N !== 0; y++) {
                for (let x = 0; x < 22 && delta_N !== 0; x++) {
                    if (starting_canvas[y][x] === 0) {
                        const dirs = [
                            [1, 0], [0, 1], [-1, 0], [0, -1]
                        ];
                        for (const [dx, dy] of dirs) {
                            if (modified_canvas[y + dy]?.[x + dx] === 1) {
                                canvas_accual.push(new Block(x, y, x + dx, y + dy, -1, -1));
                                modified_canvas[y + dy][x + dx] = 0;
                                delta_N--;
                                break;
                            }
                        }
                    }
                }
            }
        } else {
            for (let y = 0; y < 10 && delta_N !== 0; y++) {
                for (let x = 0; x < 22 && delta_N !== 0; x++) {
                    if (modified_canvas[y][x] === 0) {
                        const dirs = [
                            [1, 0], [0, 1], [-1, 0], [0, -1]
                        ];
                        for (const [dx, dy] of dirs) {
                            if (modified_canvas[y + dy]?.[x + dx] === 0) {
                                const idx = blockIndex(new Block(x, y, x, y, -1, -1), canvas_accual);
                                if (idx !== -1) {
                                    canvas_accual[idx].target_x = x + dx;
                                    canvas_accual[idx].target_y = y + dy;
                                    canvas_accual[idx].destroy_x = x + dx;
                                    canvas_accual[idx].destroy_y = y + dy;
                                    modified_canvas[y][x] = 1;
                                    delta_N++;
                                    break;
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    starting_canvas = modified_canvas.map(row => [...row]);

    voids = [];
    target_voids = [];

    for (let y = 0; y < 10; y++)
        for (let x = 0; x < 22; x++) {
            if (starting_canvas[y][x] === 1) voids.push(new Void(x, y));
            if (target_canvas[y][x] === 1) target_voids.push(new Pos(x, y));
        }

    for (let i = 0; i < voids.length; i++) {
        let min_dist = Infinity;
        let idx = -1;
        for (let j = 0; j < target_voids.length; j++) {
            const dist = Math.abs(voids[i].x - target_voids[j].x) + Math.abs(voids[i].y - target_voids[j].y);
            if (dist < min_dist) {
                min_dist = dist;
                voids[i].destination_x = target_voids[j].x;
                voids[i].destination_y = target_voids[j].y;
                idx = j;
            }
        }
        target_voids.splice(idx, 1);
    }
}

function voids_done(voids: Void[], target_canvas: number[][]): boolean {
    return !voids.some(v => target_canvas[v.y][v.x] === 0);
}

function uppdate_voids() {
    if (voids.length === 0) return;

    let idx = Math.floor(Math.random() * voids.length);
    while (voids[idx].destination_x === voids[idx].x &&
           voids[idx].destination_y === voids[idx].y &&
           voids.length > 1) {
        idx = Math.floor(Math.random() * voids.length);
    }

    const v = voids[idx];
    let dx = v.destination_x - v.x;
    let dy = v.destination_y - v.y;

    if (Math.random() > Math.abs(dx) / (Math.abs(dx) + Math.abs(dy) || 1)) {
        dx = 0; dy = Math.sign(dy);
    } else {
        dx = Math.sign(dx); dy = 0;
    }

    const bi = blockIndex(new Block(0, 0, v.x + dx, v.y + dy, 0, 0), canvas_accual);
    if (bi !== -1) {
        canvas_accual[bi].target_x = v.x;
        canvas_accual[bi].target_y = v.y;
        v.x += dx;
        v.y += dy;
    } else {
        const vi = voidIndex(new Void(v.x + dx, v.y + dy), voids);
        if (vi !== -1) {
            voids[vi].x -= dx;
            voids[vi].y -= dy;
            v.x += dx;
            v.y += dy;
        }
    }

    if (voids_done(voids, canvas_targets[target_canvas])) {
        setTimeout(() => {
            target_canvas = (target_canvas + 1) % canvas_targets.length;
            change_screen(canvas_targets[(target_canvas + canvas_targets.length - 1) % canvas_targets.length].map(r => [...r]),
                          canvas_targets[target_canvas].map(r => [...r]));
            setTimeout(() => uppdate_voids(), 500);
        }, 3000);
    } else {
        setTimeout(() => uppdate_voids(), 50);
    }
}

const canvas_targets = [[ //vlna
    [1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1],
    [1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1],
    [1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1],
    [1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0],
    [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0]
],[ //strela
    [1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1],
    [1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1],
    [1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1],
    [1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1],
    [1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1],
    [0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0]
],[ //gchd
    [1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1],
    [1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1],
    [0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0],
    [0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0],
    [0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0],
    [0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0],
    [0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0],
    [0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0],
    [1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1]
]
]


let canvas_accual: Block[] = [];
let canvas_semitarget: number[][] = [];
let target_voids: Pos[] = [];
let voids: Void[] = [];

let spacing = 29.5;
let corner = 4;
let target_canvas = 0;

function generate_random_blocks(n: number) {
    canvas_semitarget = Array.from({ length: 10 }, () => Array(22).fill(1));
    for (let i = 0; i < n; i++) {
        const x = Math.floor(Math.random() * 22);
        const y = Math.floor(Math.random() * 10);
        if (canvas_semitarget[y][x] === 1) {
            canvas_semitarget[y][x] = 0;
            canvas_accual.push(new Block(x, y, x, y, -1, -1));
        } else i--;
    }
}

$effect(() => {
    if (!canvas) return;

    loading_ctx = canvas.getContext("2d");
    if (!loading_ctx) {
        console.warn("Canvas 2D context could not be initialized.");
        return;
    }

    loading_ctx.fillStyle = "#f1effc";

    generate_random_blocks(count_blocks(canvas_targets[target_canvas]));
    draw_blocks();

    setTimeout(() => change_screen(canvas_semitarget.map(r => [...r]), canvas_targets[target_canvas].map(r => [...r])), 1000);
    setTimeout(() => uppdate_voids(), 1000);

    function animate() {
        requestAnimationFrame(animate);
        uppdate_blocks(canvas_accual, 0.3);
        draw_blocks();
    }
    animate();
});

</script>







<!-- <h1>Здравствуйте товарищи</h1> -->
<div class="grid-header blue-background">
    <div class="header1">
        <p>Inovativní matematická a fyzikální soutěž pro
        8. a 9. třídu.</p>
    </div>
    <canvas class="header2" width="660" height="300" id="this_year_canvas" bind:this={canvas}></canvas>
</div>
<div id="main">
    <div class="result-section">
        <h1>Výsledky 2025</h1>
        <div class="buttons">
            <a href="/home/results"><button class="result-button">Pražská střela online</button></a>
        </div>
    </div>
    <div class="main-a-content">
    <!-- <p style="text-align: center;">
        <b style="color: red; font-size: 18px; font-weight: bold;">Upozornění:</b> Pokud vám nepřišel email s herním odkazem, napište na strela-vlna@gchd.cz a my vám pošleme nový.
    </p> -->
    <div class="grid-main-a">
        <div class="section-title m-a0">
        <h1>Jak to probíhá</h1>
        </div>
        <p class="main-a-element blue-top m-a1">Sestavte tým ze 3 až 5 členů navštěvujících 8. a 9. třídy základní školy.</p>
        <p class="main-a-element blue-top m-a2">Zaregistrujte se do 1. 12. 2025 na Pražskou střelu (matematika) a do 3. 12. 2025 na Dopplerovu vlnu (fyzika) a 2. 12. 2025, respektive 4. 12. 2025 soutěžte!</p>
        <p class="main-a-element blue-top m-a3">Pokuste se dostat do nejlepších 15 týmů a postupte do prezenčního finále na naší škole (GCHD).</p>
        <p class="main-a-element blue-top m-a4">Prezenční kolo Dopplerovy vlny se koná 11. 12. 2025, Pražská střela pak 16. 12. 2025 na našem gymnáziu. Na ty nejlepší z vás čekají nejen hodnotné ceny, ale i teoretická možnost výhody při přijímacích zkouškách!</p>
    </div>
    <div class="video blue-top ">
        <b class="f-20">Zaujaly vás soutěže?</b> <span class="f-18">Podívejte se na naše <a href="https://www.youtube.com/watch?v=MD9si9OnKHI" class="link" target="_blank">úvodní video</a></span>
    </div>
    </div>
    <div class="main-c">
    <div class="grid-main-c">
        <h1 class="section-title m-c0">Ukázka úloh</h1>
        <div class="blue-top main-c-element m-c1">
        <h2>Pražská střela</h2>
        <h3>Lehká úloha [A]</h3>
        <p class="f-18">Pokud Kaufland zdraží 12korunovou Fidorku o 22 %, následně zdraží o 17 % z nové ceny a nakonec uvede nabídku -51 % z ceny pro uživatele Kaufland card, za kolik jsme schopni si Fidorku koupit? Cenu zaokrouhlete na celé koruny.</p>
        </div>
        <div class="blue-top main-c-element m-c2">
        <h2>Dopplerova vlna</h2>
        <h3>Lehká úloha [A]</h3>
        <p class="f-18">Martin si chtěl namíchat ideální vodu se sirupem. Proto nejdříve do malé sklenice (200ml) nalil vodu a do panákové skleničky (60ml) nalil sirup. Poté do velké sklenice nalil tři čtvrtiny malé sklenice a pět dvanáctin panákové skleničky.V jakém poměru bude sirup a voda? (Napište zlomkem v základním tvaru)</p>
        </div>
        
    </div>
    <div class="task-button">
        <div class="button-block"></div><a href="home/examples"><button class="more-tasks-button">Více úloh</button></a>
    </div>
    </div>
    <div class="grid-main-b">
    <h1 class="section-title m-b0">Harmonogram</h1>
    <div class="blue-left main-b-element m-b1">
        <h2>9:40</h2>
        <p class="f-18">Spuštění úvodního streamu na <a href="https://www.youtube.com/channel/UC8hsjufLjcGySi79GDzqJGg" class="link" target="_blank">YouTube</a></p>
    </div>
    <div class="blue-left main-b-element m-b2">
        <h2>10:00</h2>
        <p class="f-18">Začátek soutěže</p>
    </div>
    <div class="blue-left main-b-element m-b3">
        <h2>11:30</h2>
        <p class="f-18">Konec soutěže</p>
    </div>
    <div class="blue-left main-b-element m-b4">
        <h2>12:00</h2>
        <p class="f-18">Vyhlášení výsledků na streamu</p>
    </div>
    </div>
</div>


    
<style>
    .result-button{
    cursor: pointer;
    font-family: 'Lexend';
    font-size: 20px;
    font-weight: 500;
    padding: 10px;  
    padding-inline: 12px;
    color: white;
    border-radius: 10px;
    background-color: var(--color-blue);
    padding-top: 7px;
    width: 250px;
    margin: auto;
    transition: all 0.3s ease-out;
}

.buttons{
    display: flex;
    justify-content: start;
    gap: 2vw;
    padding-bottom: 20px;
    padding-top: 10px;
    margin-left: 50px;
}

.result-button:hover{
    background-color: #3c1fe2;
    box-shadow: 0px 0px 20px 0px #907cff7e;
    transition: all 0s;
}


.result-section{
      padding-left: 15%;
    padding-right: 5%;
    background-color: white;
    padding-top: 20px;
    padding-bottom: 60px;
}


    .blue-background{
        background: linear-gradient(45deg, var(--color-lightblue) 30%, var(--color-blue) 60%);
        /* background-color: var(--color-lightblue); */
    }   
.header1 {
    max-width: 700px;
    margin: 5%;
    margin-left: 10vw;
    min-width: 30vw;
    height: 300px;
    font-size: 48px;
    font-weight: 700;
    color: white;
    line-height: 60px;
    font-family: 'Lexend';
}
.header2{
    margin: auto;
    max-width: 45vw;
    padding-right: 5vw;
}

.grid-header{
    display: grid;
    max-width: 100vw;
    grid-template: 'h1 h2';
}
.link{
    font-weight: 900;
    text-decoration: none;
    color: black;
}
.link::after {
    content: '\f35d'; 
    font-family: 'Font Awesome 6 Free'; 
    font-weight: 1000; 
    font-size: 0.9em; 
    margin-left: 0.2em; 
    display: inline-block;
    vertical-align: middle;
    color: #000;
    padding-bottom: 3px;
}
#main{
    font-family: 'Lexend';
}
.f-18{
    font-size: 18px;
}
button{
    outline: none;
    border-radius: 3px;
    border: none;
    cursor: pointer;
    padding: 12px;
    font-family: 'Lexend';
}
h1{
    font-size: 40px;
}
h2{
    font-size: 25px;
    font-weight: 700;
}
h3{
    font-size: 18px;
    font-weight: 600;
}
.m-a0 {grid-area: ma0}
.m-a1 {grid-area: ma1}
.m-a2 {grid-area: ma2}
.m-a3 {grid-area: ma3}
.m-a4 {grid-area: ma4}
.m-b0 {grid-area: mb0}
.m-b1 {grid-area: mb1}
.m-b2 {grid-area: mb2}
.m-b3 {grid-area: mb3}
.m-c0 {grid-area: mc0}
.m-c1 {grid-area: mc1}
.m-c2 {grid-area: mc2}
.header1 {grid-area: h1}
.header2 {grid-area: h2}
.grid-main-a{
    display: grid;
    grid-template:
        'ma0 ma0 ma0 ma0'
        'ma1 ma2 ma3 ma4'; 
        
}

.main-a-element{
    font-size: 18px;
    max-width: 400px;
    margin-inline: 10%;
}

.grid-main-b{
    padding-left: 15%;
    max-width: 100vw;
    padding-right: 5%;
    background-color: var(--color-light-gray);
    padding-top: 20px;
    padding-bottom: 60px;
    display: grid;
    grid-template: 
    'mb0'
    'mb1'
    'mb2'
    'mb3'
}

.grid-main-c{  
    display: grid;
    padding-left: 15%;
    max-width: 100vw;
    padding-right: 5%;
    padding-top: 40px;
    grid-template: 
        'mc0 mc0'
        'mc1 mc2';
    justify-content: start;
    column-gap: 10vw;
}

.main-c{
    background-color: white
}
.main-c-element{
    max-width: 500px;
    margin-left: 30px;
}

.main-b-element{
    margin-top: 10px;
    margin-bottom: 10px;
    max-width: 400px;
    margin-left: 30px;
}
.more-tasks-button{
    cursor: pointer;
    font-family: 'Lexend';
    font-size: 23px;
    font-weight: 500;
    padding: 13px;  
    padding-inline: 15px;
    border-radius: 10px;
    color: white;
    background-color: var(--color-blue);
    padding-top: 10px;
    width: 200px;
    margin: auto;
    transition: all 0.3s ease-out;
}

.task-button{
    padding-left: 15%;
    padding-right: 5%;
    display: flex;
    justify-content: start;
    padding-bottom: 50px;
    padding-top: 50px;
}

.button-block{
    width: 38%;
    max-width: 530px;
}
.more-tasks-button:hover{
    background-color: #3c1fe2;
    box-shadow: 0px 0px 20px 0px #907cff7e;
    transition: all 0s;
}

.main-a-content{
    display: flex;
    flex-direction: column;
    padding-left: 15%;
    padding-right: 5%;
    background-color: var(--color-light-gray);
    padding-top: 20px;
    padding-bottom: 60px;
    gap: 30px;
}

.video{
    padding-top: 20px;
    padding-bottom: 20px;
    width: 90%;
    margin: auto;
    display: flex;
    justify-content: space-between;
}

.video span{
    margin-right: 2vw;
}

.video b{
    margin-left: 2vw;
}
.blue-top{
    border-top: 4px var(--color-blue) solid;
    border-radius: 4px;
    padding-top: 10px;
}
.blue-left {
    border-left: 4px var(--color-blue) solid;
    border-radius: 4px;
    padding-left: 10px;
}
@media (max-width: 1100px){
    .grid-main-c, .task-button{
        padding-left: 5%;
    }
    .grid-main-a{
        margin-left: 5%;
    }
    .main-a-content{
        padding-left: 5%;
    }
    .header1{
    margin: auto;
    }
    .header2{
        margin: auto;
        max-width: 90vw;
        padding: 0px;
    }
    .grid-header{
        grid-template: 
            'h1'
            'h2'
    }
    
}
@media (max-width: 700px){
    .grid-main-a{
        grid-template: 
        'ma0'
        'ma1'
        'ma2'
        'ma3'
        'ma4'
    }
    .main-a-content{
        padding-left: 5%;
        gap: 50px;
    }
    .grid-main-c{
        grid-template: 
        'mc0'
        'mc1'
        'mc2'
        'mc3'
    }
    .grid-main-b{
        padding-left: 5%;
    }
    .main-b-element{
        max-width: 100vw;
    }
    .button-block{
        display: none;
    }
    .task-button{
        justify-content: center;
    }
    .m-c2{
        margin-top: 60px;
    }
    .video{
        flex-direction: column;
        gap: 20px;
        align-items: center;
        max-width: 400px;
    }
    .grid-header{
    grid-template: 
        'h1'
        'h2'
    } 
    .header2{
        max-width: 80vw;
    }      
    .header1{
        min-width: 80%;
        font-size: 40px;
        line-height: 55px;
        height: auto;
        margin: auto;
        max-width: 80vw;
    }
}
@media (max-width:300px){
    .m-b0{
        font-size: 12vw
    }
    .header1{
    font-size: 14vw;
    line-height: 17vw;
    }
}

.section-title{
    margin-top: 10px;
}
</style>
