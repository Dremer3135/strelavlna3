<script lang="ts">
    import LoadingAnimation from "../../../../../shared/dist/components/LoadingAnimation.svelte";

  let { isPaused }: { isPaused: boolean } = $props();

</script>

{#if isPaused}
  <div class="paused-overlay">
    <div class="paused-content">
      <div class="title">
        {#each "Soutěž je pozastavena z technických důvodů...".split(" ") as word, index}
          <h1 class="chunk" style="animation-delay: {(index * 0.1)}s; color: color-mix(in srgb, var(--color-{["pink", "orange", "yellow", "purple"][index%4]}) 90%, black 10%);">{word}</h1>
        {/each}
        <!-- Soutěž je momentálně pozastavena z technických důvodů... -->
      </div>
      <p>Prosím, vyčkejte na obnovení soutěže.</p>
      <p>Na řešení problému pracujeme.</p>
      <div class="animation">
        <LoadingAnimation color="black"/>
      </div>
    </div>
  </div>
{/if}

<style lang="scss">
  .paused-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    backdrop-filter: blur(4px);
    /* background-color: rgba(0, 0, 0, 0.5); */
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 9999;
  }

  .paused-content {
    font-family: 'Lexend';
    background-color: white;
    padding: 40px;
    border-radius: 8px;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;

    p {
      margin: 3px;
    }

    .title {
      display: flex;
      flex-direction: row;
      justify-content: center;
      gap: 10px;
      
      h1 {
        // transform: translateY(100px);
        animation: animation 4s cubic-bezier(0.215, 0.610, 0.355, 1) infinite;
        color: color-mix(in srgb, var(--color-purple) 30%, black 70%);
      }

    }
  }

  .animation {
    margin-top: 150px;
    width: 100px;
    height: 100px;
  }

  @keyframes animation {
    0% { transform: translateY(0px); }
    15% { transform: translateY(0px); }
    30% { transform: translateY(-20px); }
    60% { transform: translateY(-20px); }
    75% { transform: translateY(0px); }
    100% { transform: translateY(0px); }
  }


</style>