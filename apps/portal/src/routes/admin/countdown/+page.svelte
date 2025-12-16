<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    // import skrat from "$lib/assets/images/SKRAT.svg";

    const DT = 1/60;

    const SPAWN_DIST = 100
    const EXPLOSION_OFFSET = 10; // Distance in front of camera to explode

    const FDIST = 1000;
    const SIZE = 0.5;
    const CAMERA_SPEED = 5;

    const collors: string[] = ["#9500EB", "#EB0072", "#EB6E00", "#EBAD00"];

    let width = $state(0);
    let height = $state(0);

    let canvas: HTMLCanvasElement;
    let ctx: CanvasRenderingContext2D | null;

    let xOftset = 0;
    let xOftsetTarget = 0;

    let showDatePicker = $state(false);
    
    // Helper to format Date to "YYYY-MM-DDTHH:mm" for input
    const formatDateForInput = (date: Date) => {
        const offset = date.getTimezoneOffset();
        const localDate = new Date(date.getTime() - (offset*60*1000));
        return localDate.toISOString().slice(0, 16);
    };

    let startDateVal = $state(formatDateForInput(new Date()));
    let endDateVal = $state(formatDateForInput(new Date(new Date().getTime() + 60 * 60 * 1000))); // Default +1 hour
    let fontSizeVal = $state(100);

    let timeLeft = $state("00:00:00");
    let blinking = $state(false);

    let creditsHovered = $state(false);

    let interval: any;
    let animationFrameId: number;

    // Time sync state
    let startTime = 0;
    let nextBeatIndex = 0; // Tracks which second we need to spawn next

    onMount(() => {
        ctx = canvas.getContext("2d");
        
        // Initialize start time for sync
        startTime = Date.now();
        // Determine the first beat index based on current time
        nextBeatIndex = Math.ceil((Date.now() - startTime) / 1000);

        animationFrame();

        interval = setInterval(updateTime, 1000);
        updateTime();
    });

    onDestroy(() => {
        if (interval) clearInterval(interval);
        if (animationFrameId) cancelAnimationFrame(animationFrameId);
    });

    function updateTime() {
        const now = new Date();
        const start = new Date(startDateVal);
        const end = new Date(endDateVal);
        let diff = 0;

        blinking = false;

        if (now < start) {
            diff = start.getTime() - now.getTime();
        } else if (now < end) {
            diff = end.getTime() - now.getTime();
        } else {
            timeLeft = "--:--:--";
            blinking = true;
            return;
        }

        const hours = Math.floor(diff / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);

        timeLeft = `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    
    type Vector2 = {
        x: number,
        y: number
    }
    
    type Vector3 = {
        x: number,
        y: number,
        z: number
    }

    type FallingCube = {
        position: Vector3,
        // speed: Vector2,
        orientation: number,
        angularSpeed: number,
        color: number,
        type: 'random' | 'beat',
        // size: number,
        // paralaxCof: number
    }

    type Particle = {
        position: Vector3,
        velocity: Vector3,
        life: number,
        maxLife: number,
        color: string
    }

    let fallingCubes: FallingCube[] = [];
    let particles: Particle[] = [];
    let cameraPosition: Vector3 = {x: 0, y: 0, z: 0}

    function addRandomCube() {
        // Only spawn if within the time window
        const now = Date.now();
        const start = new Date(startDateVal).getTime();
        const end = new Date(endDateVal).getTime();

        // Shift window by flight time (20s) so cubes appear at camera at start time
        const flightTime = (SPAWN_DIST / CAMERA_SPEED) * 1000;
        if (now < start - flightTime || now > end - flightTime) return;

        // Spawn random cubes ahead of the camera
        let x = 0, y = 0;
        
        // Ensure we don't spawn in the central corridor (where beat cubes are)
        do {
            x = (Math.random() - 0.5) * 50;
            y = (Math.random() - 0.5) * 50;
        } while (Math.abs(x) < 8 && Math.abs(y) < 8); // Exclusion zone slightly larger than beat area

        let position: Vector3 = { x, y, z: cameraPosition.z + SPAWN_DIST };
        let orientation: number = Math.random()*360;
        let angularSpeed: number = (Math.random()*2 - 1) * 90;
        let color: number = Math.floor(Math.random()*4);
        // fallingCubes.push();
        fallingCubes = [{ position, orientation, angularSpeed, color, type: 'random' }, ...fallingCubes];
    }

    function spawnBeatCubes() {
        // Calculate the Z horizon we want to populate
        // We want to ensure cubes are spawned up to SPAWN_DIST ahead
        const horizonZ = cameraPosition.z + SPAWN_DIST;
        
        // Check if the next beat index is within the horizon
        while (true) {
            // Calculate absolute time for this beat
            const beatTime = startTime + nextBeatIndex * 1000;
            const beatZ = (beatTime - startTime) / 1000 * CAMERA_SPEED + EXPLOSION_OFFSET;

            // If this beat is beyond the horizon, stop spawning
            if (beatZ > horizonZ) break;

            // Only spawn if within the start/end window
            const start = new Date(startDateVal).getTime();
            const end = new Date(endDateVal).getTime();

            // Add small buffer to end time to ensure last cube isn't skipped due to float math
            if (beatTime >= start && beatTime <= end + 100) {
                let position: Vector3 = { 
                    x: (Math.random() - 0.5) * 10, // Random X within central area (tighter)
                    y: (Math.random() - 0.5) * 10, // Random Y within central area (tighter)
                    z: beatZ
                };
                let orientation: number = Math.random() * 360; 
                let angularSpeed: number = (Math.random() * 2 - 1) * 90;
                let color: number = nextBeatIndex % 4; // Cycle colors
                
                fallingCubes = [{ position, orientation, angularSpeed, color, type: 'beat' }, ...fallingCubes];
            }
            
            nextBeatIndex++;
        }
    }

    function createExplosion(pos: Vector3, colorIdx: number) {
        const color = collors[colorIdx];
        for (let i = 0; i < 40; i++) {
            const speed = 5 + Math.random() * 10; // Faster particles
            const angle = Math.random() * Math.PI * 2;
            const zSpeed = (Math.random() - 0.5) * 30; // High Z speed for 3D effect
            
            particles.push({
                position: { x: pos.x, y: pos.y, z: pos.z }, 
                velocity: {
                    x: Math.cos(angle) * speed,
                    y: Math.sin(angle) * speed,
                    z: zSpeed
                },
                life: 1.0,
                maxLife: 1.0,
                color: color
            });
        }
    }

    function updateParticles(dt: number) {
        for (let p of particles) {
            p.position.x += p.velocity.x * dt;
            p.position.y += p.velocity.y * dt;
            p.position.z += p.velocity.z * dt;
            p.life -= dt * 1.5; // Fade out slightly faster
        }
        particles = particles.filter(p => p.life > 0);
    }

    function updateCubes(dt: number) {
        let cubesToRemove: number[] = [];

        // Update orientations
        for (let cube of fallingCubes) {
            cube.orientation += cube.angularSpeed * dt;
        }

        // Check for beat cubes passing camera
        // Since we modify the array, we iterate carefully or filter
        let nextCubes: FallingCube[] = [];
        
        for (let cube of fallingCubes) {
            // Logic for beat cubes explosion
            if (cube.type === 'beat') {
                if (cube.position.z <= cameraPosition.z + EXPLOSION_OFFSET) {
                    createExplosion(cube.position, cube.color);
                    continue; // Remove this cube (explode)
                }
            } else {
                // Random cubes just despawn when behind camera
                if (cube.position.z <= cameraPosition.z) {
                    continue;
                }
            }
            nextCubes.push(cube);
        }
        fallingCubes = nextCubes;
    }

    function renderParticles() {
        if (!ctx) return;
        
        let screenCenter: Vector2 = {x: width/2, y: height/2};

        for (let p of particles) {
            let dx = p.position.x - cameraPosition.x;
            let dy = p.position.y - cameraPosition.y;
            let dz = p.position.z - cameraPosition.z;

            // Simple projection - if particle is behind camera, don't draw (unless very close)
            // Allow drawing slightly behind for explosion effect feel? No, standard projection.
            if (dz < 0.1) continue; // Don't draw if behind or too close

            let x = screenCenter.x + dx * (FDIST/dz);
            let y = screenCenter.y + dy * (FDIST/dz);

            let size = (FDIST / dz * SIZE) * 0.25; // Particles size (half radius)

            ctx.fillStyle = p.color;
            ctx.globalAlpha = Math.max(0, p.life);
            
            ctx.beginPath();
            ctx.arc(x, y, size, 0, Math.PI * 2);
            ctx.fill();
        }
        ctx.globalAlpha = 1;
    }

    function renderCubes() {
        if(!ctx) return;

        // console.log("heheheha");

        ctx.clearRect(0, 0, width, height);
        
        let screenCenter: Vector2 = {x: width/2, y: height/2};
        // console.log(screenCenter.y);

        for(let cube of fallingCubes) {
            let dx = cube.position.x - cameraPosition.x;
            let dy = cube.position.y - cameraPosition.y;
            let dz = cube.position.z - cameraPosition.z;

            let x = screenCenter.x + dx * (FDIST/dz);
            let y = screenCenter.y + dy * (FDIST/dz);

            let opacity = 1;
            
            if (cube.type === 'random') {
                opacity = Math.max(0, Math.min(1, (SPAWN_DIST - dz) / 20));
            } else {
                // Beat cubes shouldn't fade in too aggressively, or maybe they should?
                // Let's keep them visible.
                opacity = Math.max(0, Math.min(1, (SPAWN_DIST - dz) / 10));
            }
            
            if (opacity <= 0) continue;

            let size = FDIST / dz * SIZE;
            
            // console.log(x, y);

            ctx.fillStyle = collors[cube.color];
            ctx.save();
            ctx.globalAlpha = opacity;
            ctx.translate(x, y);
            ctx.rotate(cube.orientation * Math.PI / 180);

            ctx.beginPath();
            ctx.roundRect(-size/2, -size/2, size, size, 6);
            ctx.fill();
            
            ctx.restore();
        }
    }

    function updateParalax() {
        xOftset += (xOftsetTarget - xOftset) * 0.05;
    }

    function updateCamera() {
        if (startTime === 0) return;
        // Sync camera to absolute time
        let now = Date.now();
        let elapsed = (now - startTime) / 1000;
        cameraPosition.z = elapsed * CAMERA_SPEED;
    }


    let lastTime = 0;

    function animationFrame() {
        animationFrameId = requestAnimationFrame(animationFrame);
        
        let now = performance.now() / 1000;
        if (lastTime === 0) lastTime = now;
        let dt = now - lastTime;
        lastTime = now;
        
        // Cap dt to avoid huge jumps if tab inactive
        if (dt > 0.1) dt = 0.1;

        if (Math.random() < 0.2) {
            addRandomCube();
        }

        // updateParalax();
        updateCamera();
        spawnBeatCubes();
        updateCubes(dt);
        updateParticles(dt);
        
        renderCubes();
        renderParticles();

        // console.log("lalal");
    }


</script>

<main bind:clientWidth={width} bind:clientHeight={height} onmousemove={(e) => {
    xOftsetTarget = (e.clientX - width) / 5;
}}>
    <div class="title" onclick={() => {showDatePicker = !showDatePicker}}>
        <h1 class:blinking={blinking} style="font-size: {fontSizeVal}px;">{timeLeft}</h1>
        <!-- {#if type == "before"}
            <h1>Vyčkejte na začátek soutěže</h1>
        {:else if type == "waiting-for-results"}
            <h1>Právě vyhodnocujeme výsledky</h1>
        {:else}
            <h1>{$currentState.teamName}</h1>
        {/if} -->
        {#if showDatePicker}
            <div class="inputs" onclick={(e) => e.stopPropagation()}>
                <label>
                    Start:
                    <input type="datetime-local" bind:value={startDateVal}>
                </label>
                <label>
                    End:
                    <input type="datetime-local" bind:value={endDateVal}>
                </label>
                <label>
                    FontSize:
                    <input type="number" bind:value={fontSizeVal}>
                </label>
            </div>
        {/if}
    </div>
    <!-- {#if type=="results"}
        <div class="results">
            <div class="stats">
                <div class="rank property-wrapper">
                    <h2 class="value">{$currentState.rank}.</h2>
                    <p class="name">Umístění</p>
                </div>
                <div class="money property-wrapper">
                    <h2 class="value">{$currentState.money}</h2>
                    <p class="name">Body</p>
                </div>
            </div>
            {#if $currentState.rank > 15}
                <h2 class="message">Pokud jste se umístili v top 15, dostanete potvrzující email s dalšími instrukcemi o postupu do finálového kola</h2>
            {:else}
                <h2 class="message">Vypadá to, že jste se umístili v top 15! Očekávejte tedy potvrzující email s dalšími instrukcemi o postupu do finálového kola</h2>
            {/if}
        </div>
    {/if} -->
    <div class="credits">
        <!-- <img src={skrat} alt="credits to skrat team" onmouseenter={() => { creditsHovered = true; }} onmouseleave={() => { creditsHovered = false; }}> -->
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

        canvas {
            z-index: -1;
            background-color: black;
            filter: invert(0%);
        }

        .title {
            position: absolute;
            top: 0px;
            left: 0px;
            width: 100%;
            height: 100%;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;

            h1 {
                font-family: 'Fredoka';
                margin: 0px;

                letter-spacing: 5px;
                font-size: 100px;
                font-weight: 600;
                box-sizing: border-box;
                border-radius: 5px;
                padding: 9px 25px;
                color: color-mix(in srgb, var(--color-purple) 40%, black 60%);
                color: white;
                filter: drop-shadow(3px -10px 0px #ffffff50);

                &.blinking {
                    animation: blink 1s infinite;
                }
            }

            .inputs {
                margin-top: 20px;
                display: flex;
                flex-direction: row;
                gap: 20px;
                background-color: #F0F0F030;
                backdrop-filter: blur(3px);
                padding: 15px;
                border-radius: 5px;
                
                label {
                    display: flex;
                    flex-direction: column;
                    color: white;
                    font-family: 'Fredoka';
                    font-weight: 500;
                }

                input {
                    margin-top: 5px;
                    padding: 5px;
                    border-radius: 3px;
                    border: none;
                }
            }

        }

        @keyframes blink {
            0% { opacity: 1; }
            50% { opacity: 0; }
            100% { opacity: 1; }
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

        .results {
            position: absolute;
            top: 200px;

            .stats {
                position: relative;
                width: 100%;
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: space-evenly;
    
                .property-wrapper {
                    display: flex;
                    flex-direction: column;
                    gap: 15px;
                    justify-content: center;
                    align-items: center;
                    padding-top: 100px;
                    backdrop-filter: blur(3px);
                    background-color: #F0F0F030;
                    box-sizing: border-box;
                    box-shadow: 0px 0px 10px #F0F0F0;
                    border-radius: 5px;
                    padding: 9px 25px;
    
                    .value {
                        font-family: 'Lexend';
                        font-size: 60px;
                        margin: 0px;
                        padding: 10px 60px;
    
                        border-bottom: 3px color-mix(in srgb, var(--color-purple) 30%, black 70%) dashed;
                    }
    
                    .name {
                        font-family: 'Fredoka';
                        font-size: 25px;
                        color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
                        font-weight: 500;
                        margin: 0px;
                    }
                }
    
    
    
                .money {
                    .value {
                        color: var(--color-pink);
                    }
                }
                .rank {
                    .value {
                        color: var(--color-yellow);
                    }
                }
            }
    
            .message {
                margin-top: 70px;
                font-family: 'Fredoka';
                font-weight: 500;
                color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
                backdrop-filter: blur(3px);
                background-color: #F0F0F030;
                box-sizing: border-box;
                box-shadow: 0px 0px 10px #F0F0F0;
                border-radius: 5px;
                padding: 9px 25px;
            }
        }
    }
</style>