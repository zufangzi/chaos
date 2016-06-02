$(function() {
    
    $("#mainReturn").click(
        function() {
            $("#cloudHeader li a").removeAttr("style")
             $("#cloudHeader li").removeAttr("style")
        }
    );

    $("#cloudHeader li a").click(
        function() {
            $(this).css("color", "#fff")
            $(this).parent('li').css("background", "rgba(255, 255, 255, 0.2)")
            $(this).parent('li').siblings().removeAttr("style")
            $(this).parent('li').siblings().find('a').removeAttr("style")
        }
    );

    $('.banner').unslider({
        speed: 500,
        delay: 3000,
        keys: true,
        dots: true,
        fluid: true,
        autoplay: true,
        arrows: false

    });
});