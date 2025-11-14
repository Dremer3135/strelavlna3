<script lang="ts">
  import { createEventDispatcher } from "svelte";
  let { team } = $props()
  // export let team: any;

  const dispatch = createEventDispatcher();
</script>

<div class="modal-backdrop" on:click={() => dispatch("close")}>
  <div class="modal-content" on:click|stopPropagation>
    <h2>{team.name}</h2>
    <ul>
      {#if team.player1email}<li>{team.player1email}</li>{/if}
      {#if team.player2email}<li>{team.player2email}</li>{/if}
      {#if team.player3email}<li>{team.player3email}</li>{/if}
      {#if team.player4email}<li>{team.player4email}</li>{/if}
      {#if team.player5email}<li>{team.player5email}</li>{/if}
    </ul>
    <button on:click={() => dispatch("close")}>Close</button>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 102;
  }
  .modal-content {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    color: #333;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    padding: 0.5rem 0;
  }
</style>
