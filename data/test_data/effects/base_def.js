
var ability_scores = ["str","dex","con","int","wis","cha"];

var bonus_types = ["alchemical", "armor", "armor_enhancement", "circumstance",
    "competence", "deflection", "dodge", "enhancement", "inherent", "insight",
    "luck", "morale", "natural_armor", "natural_armor_enhancement", "profane",
    "racial", "resistance", "sacred", "shield", "shield_enhancement", "size",
    "trait", "untyped"];
var touch_ac_exclude = ["armor", "armor_enhancement", "natural_armor",
    "natural_armor_enhancement", "shield", "shield_enhancement"];
var flatfooted_ac_exclude = ["dodge"];

var str_skills = ["skill_climb", "skill_swim"];
var dex_skills = ["skill_acrobatics", "skill_disabledevice",
    "skill_escapeartist", "skill_fly", "skill_ride", "skill_sleightofhand",
    "skill_stealth"];
var con_skills = [];
var int_skills = ["skill_appraise", "skill_craft", "skill_knowledge",
    "skill_linguistics", "skill_spellcraft"];
var wis_skills = ["skill_heal", "skill_perception", "skill_profession",
    "skill_sensemotive", "skill_survival"];
var cha_skills = ["skill_bluff", "skill_diplomacy", "skill_disguise",
    "skill_handleanimal", "skill_intimidate", "skill_perform",
    "skill_usemagicdevice"];

define.rule({
  id: "baseStatsModifiers",
  eval: function(entity) {
    entity.vars.new({
      id: "str_mod",
      init: 0,
    });
    entity.reductions.new({
      id: "baseStatsModifiers",
      depends: ability_scores.map(function(stat) { return stat+"_base";}),
      modifies: ["str_mod"],
      eval: function(deps, mods) {
        mods.str_mod.add(Math.round((deps.str_base - 10)/2));
      },
    })
  },
})

define.effect({
  id: "baseEntityRules",
  displayName: "Base Rules",
  displayType: "yep",
  rules: ["baseStatsModifiers"],
})
