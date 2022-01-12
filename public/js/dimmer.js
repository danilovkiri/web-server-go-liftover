function dim(bool)
{
    if (typeof bool=='undefined') bool=true;
    document.getElementById('dimmer').style.display=(bool?'block':'none');
}
