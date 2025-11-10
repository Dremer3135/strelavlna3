<script lang="ts">
  import { pocketbase } from "$lib/pocketbase";
  import Navbar from "$lib/components/register/Navbar.svelte";
  import { contestsStore, fetchTeams } from "$lib/stores/register";
  import { currentUser } from '$lib/stores/auth';
  import Contest_simple from "$lib/components/register/Contest_simple.svelte";
  import Auth from "$lib/components/register/Auth.svelte";
  import TeamsPanel from "$lib/components/register/TeamsPanel.svelte";
  import AddTeamModal from "$lib/components/register/AddTeamModal.svelte";
  import TeamDetailsModal from "$lib/components/register/TeamDetailsModal.svelte";
  import blocks_clumped from "$lib/assets/images/register/blocks_clumped.svg?url";

  let showAuth = $state(false);
  let authType = $state<"login" | "register">("login");
  let showTeamsPanel = $state(false);
  let showAddTeam = $state(false);
  let showTeamDetails = $state(false);
  let selectedTeam = $state<any>(null);

  $effect(() => {
      if ($currentUser) {
          fetchTeams($currentUser.id);
      }
  });

</script>

<div id="app">
  <Navbar currentUser={$currentUser} on:showAuth={(e) => {console.log('showAuth detail:', e.detail); showAuth = true; authType = e.detail}} on:toggleTeamsPanel={() => (showTeamsPanel = !showTeamsPanel)} />
  
  {#if showAuth}
    <Auth on:close={() => (showAuth = false)} type={authType} />
  {/if}
<!-- 
  {#if showTeamsPanel}
    <TeamsPanel 
      on:close={() => (showTeamsPanel = false)} 
      on:showAddTeam={() => (showAddTeam = true)}
      on:showTeamDetails={(e) => {selectedTeam = e.detail; showTeamDetails = true;}}
    />
  {/if} -->
<!-- 
  {#if showAddTeam}
      <AddTeamModal on:close={() => (showAddTeam = false)} />
  {/if} -->

  {#if showTeamDetails}
      <TeamDetailsModal team={selectedTeam} on:close={() => (showTeamDetails = false)} />
  {/if}

  <div class="content-holder">
    <div class="logout-page">
      <h1 class="upcomming"><span class="anchor"><img src={blocks_clumped} class="blocks-clumped"></span> Nadcházející</h1>
      <div class="contests-list">
        {#each $contestsStore as contest}
          <Contest_simple {contest} />
        {/each}
      </div>
    </div>
  </div>
</div>

<style>
  .contests-list {
    margin-top: 2rem;
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 100px;
  }

  .logout-page {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 80px;
    max-width: 1300px; 
    width: 100%;
    box-sizing: border-box;
  }

  .upcomming {
    font-family: "Fredoka";
    font-size: 50px;
    font-weight: bold;
    color: #002C5E;

    margin-bottom: 100px;
  }

  .content-holder {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .anchor {
    position: relative;
  }

  .blocks-clumped {
    position: absolute;
    top: 55px;
    left: -20px;
  }
</style>
