
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
    entity.vars.new({
      id: "fighter_progressions_proxy",
      dependencies: ["fighter_lvl"],
      modifies: ["bab", "will_save", "fort_save", "ref_save"],

      onEval: function(deps, mods) {
        mods.bab(deps.fighter_lvl);
        mods.will_save(Math.floor(deps.fighter_lvl/3));
        mods.ref_save(Math.floor(deps.fighter_lvl/3));
        mods.fort_save(Math.floor(deps.fighter_lvl/2)+2);
      }
    });
  }
})
