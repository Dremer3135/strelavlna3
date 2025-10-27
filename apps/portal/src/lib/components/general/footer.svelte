<script>
    import yt_icon from '$lib/assets/images/icons/youtube_icon.png';
    import in_icon from '$lib/assets/images/icons/instagram_icon.png';
    import gh_icon from '$lib/assets/images/icons/github_icon.png';
    import skrat_icon from '$lib/assets/images/icons/skrat.png';
</script>

<footer class="footer">
    <div class="grid-footer">
        <div class="wrapper">
            <div class="footer-element1 footer-element">
                <h2>Kde nás najdete?</h2>
                <div class="gray-top">
                    <p><b>Gymnázium Christana Dopplera</b></p>
                    <p>Zborovská 45, Praha 5 - Malá Strana 150 00</p>
                    <p>Kontaktní email: <b>strela-vlna@gchd.cz</b></p>
                </div>
            </div>
        </div>
        <div class="wrapper">
            <div class="footer-element2 footer-element">
                <h2>Mapa</h2>
                <div class="gray-top">
                    <p><iframe style="border:none" title="mapa" src="https://frame.mapy.cz/s/podetabahu" height="200" width="100%" frameborder="0"></iframe></p>
                </div>
            </div>
        </div>
        <div class="wrapper">
            <div class="footer-element3 footer-element">
                <h2>Doprava MHD</h2>
                <div class="gray-top">
                    <p>Tramvajová zastávka <b>Újezd</b> (tramvaje 9, 12, 15, 20, 22)</p>
                    <p>Metro B - <b>Národní třída</b> a poté tramvají 9 nebo 22</p>
                </div>
            </div>
        </div>
        <div class="wrapper">
            <div class="footer-element4 footer-element" >
                <h2>Další odkazy</h2>
                <div class="gray-top">
                    <p>
                        <a href="https://www.youtube.com/channel/UC8hsjufLjcGySi79GDzqJGg" title="YouTube" target="_blank"><img alt="YouTube_icon" src="{yt_icon}" width="48" height="34"></a>
                        <a href="https://www.instagram.com/strela_vlna?igsh=MXVkMzkyMXdjajBwdg==" title="Instagram" target="_blank"><img alt="Instagram_icon" src="{in_icon}" width="40" height="40" style="padding-left: 15px"></a>
                        <a href="https://github.com/anteat3r/strelavlna2" title="GitHub" target="_blank"><img alt="GitHub_icon" src="{gh_icon}" width="40" height="40" style="padding-left: 15px"></a>
                    </p>
                </div>
            </div>
        </div>
    </div>
    <div class="copyright">
        <p>© 2024 Skrat</p> <img alt="skrat-logo" src = "{skrat_icon}" width="70" height="70">
    </div>
</footer>

<style>
.footer-element1 {grid-area: f1}
.footer-element2 {grid-area: f2}
.footer-element3 {grid-area: f3}
.footer-element4 {grid-area: f4}
.grid-footer{
    display: grid;
    max-width: 100vw;
    grid-template: 
        'f1 f2 f3 f4';
}
.footer-element{
    /* margin-inline: 10%; */
    display: inline-block;
    max-width: 300px;
    min-width: 160px;
    text-align: left;
    width: 20vw;
}
.wrapper{
    display: flex;
    justify-content: center;
}
.footer{
    background-color: rgb(0, 0, 0);
    color: white;
    padding-top: 40px;
    width: 100%;
    font-family: 'Lexend';
    margin-bottom: 0;
    margin-top: auto;
}
.gray-top{
    border-top: 3px var(--color-gray) solid;
    border-radius: 3px;
    padding-top: 5px;
}
.copyright{
    display: flex;
    align-items: center;
    justify-content: center;
    padding-top: 20px;
    padding-bottom: 35px;
    text-align: center;
    margin: auto;
}
@media (max-width: 1100px){
    .grid-footer{
        grid-template: 
            'f1 f2'
            'f3 f4';
    }
    .footer-element{
        width: 40vw;
        margin-inline: auto;
    }
}
@media (max-width: 700px){
    .grid-footer{
        grid-template: 
            'f1' 
            'f2'
            'f3'
            'f4';
    }
    .footer-element{
        margin: auto;
        width: 70vw;
    }
}
</style>