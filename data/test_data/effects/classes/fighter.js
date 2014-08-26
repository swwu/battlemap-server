
define.effect({
  id: "fighterClass",
  displayName: "Fighter",
  displayType: "class",

  onEffect: function(entity) {
    // fighter level, for fighter class features
    entity.vars.newAccum({
      id: "fighter_lvl",
      op: "+",
      init: 0,
    });
    entity.vars.newProxy({
      id: "fighter_progressions_proxy",
      depends: ["fighter_lvl"],
      modifies: ["hp", "bab", "will_save_base", "fort_save_base", "ref_save_base"],

      onEval: function(deps, mods) {
        // HD progression, using average
        mods.hp(Math.floor(deps.fighter_lvl*5.5+4.5));
        // bab, saves
        mods.bab(deps.fighter_lvl);
        mods.will_save_base(Math.floor(deps.fighter_lvl/3));
        mods.ref_save_base(Math.floor(deps.fighter_lvl/3));
        mods.fort_save_base(Math.floor(deps.fighter_lvl/2)+2);
      }
    });

    /*
    entity.labels.new({
      displayName: "Bravery",
      displayDesc: "Starting at 2nd level, a fighter gains a +1 bonus on Will saves against fear. This bonus increases by +1 for every four levels beyond 2nd.",
      displayType: "ability",

      condition: function(deps) {
      }
    });
    */

  }
})
