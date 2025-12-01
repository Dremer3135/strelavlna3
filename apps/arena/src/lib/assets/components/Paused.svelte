<script lang="ts">
    import LoadingAnimation from "../../../../../shared/dist/components/LoadingAnimation.svelte";

  let { isPaused }: { isPaused: boolean } = $props();

</script>

{#if isPaused}
  <div class="paused-overlay">
    <div class="paused-content">
      <h1>Soutěž je momentálně pozastavena z technických důvodů...</h1>
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
  }

  .animation {
    margin-top: 150px;
    width: 100px;
    height: 100px;
  }


</style>