<script lang="ts">
  import * as m from "$lib/paraglide/messages"
  import { resolve } from '$app/paths';
</script>

<footer class="footer sm:footer-horizontal footer-center bg-base-300 text-base-content p-4 flex items-center justify-center w-full">
  <aside class="flex items-center justify-center flex-row">
    <p>{ m.footer_copyright()} © {new Date().getFullYear()} - {m["footer_copyright_all-right-reserved-by"]()}</p>
    <a class="link" target="_blank" href={resolve("/terms")}>Terms of Service</a>
    <a class="link" target="_blank" href={resolve("/privacy")}>Privacy Policy</a>
  </aside>
</footer>
