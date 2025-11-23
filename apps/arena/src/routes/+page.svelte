<script lang="ts">
  import { page } from '$app/stores';
  const socket = new WebSocket(`https://sv.skrat.org/ws?token=${$page.data.token}`)
  socket.addEventListener("message", (event) => {
    console.log(typeof event.data)
  })
</script>
