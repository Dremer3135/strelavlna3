<script lang="ts">
  import { teamsStore } from "$lib/stores/register";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();
</script>

<div class="teams-panel">
  <div class="teams-panel-header">
    <h2>Moje týmy</h2>
    <button class="close-btn" on:click={() => dispatch("close")}>X</button>
  </div>
  <div class="teams-list">
    {#each $teamsStore as team}
      <div class="team-item" on:click={() => dispatch("showTeamDetails", team)}>
        {team.name}
      </div>
    {/each}
  </div>
  <button class="add-team-btn" on:click={() => dispatch("showAddTeam")}>+</button>
</div>

<style>

  .teams-panel {
    position: fixed;
    top: 0;
    right: 0;
    width: 300px;
    height: 100%;
    background-color: white;
    color: black;
    padding: 1rem;
    /* box-shadow: -2px 0 5px rgba(0, 0, 0, 0.5); */
    z-index: 101;
    display: flex;
    flex-direction: column;
    font-family: "Lexend";
  }
  .teams-panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid white;
    padding-bottom: 0.5rem;
  }
  .close-btn {
    background: none;
    border: none;
    color: white;
    font-size: 1.5rem;
    cursor: pointer;
  }
  .teams-list {
    flex-grow: 1;
    margin-top: 1rem;
  }
  .team-item {
    padding: 0.5rem;
    cursor: pointer;
    border-bottom: 1px solid #336b9f;
  }
  .team-item:hover {
    background-color: #336b9f;
  }
  .add-team-btn {
    background-color: white;
    color: #003f88;
    border: none;
    border-radius: 50%;
    width: 50px;
    height: 50px;
    font-size: 2rem;
    cursor: pointer;
    align-self: center;
    margin-top: 1rem;
  }
</style>
