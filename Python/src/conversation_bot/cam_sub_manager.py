import logging
import re

from telegram import Update, ReplyKeyboardMarkup, ReplyKeyboardRemove
from telegram.ext import (CallbackContext, ConversationHandler, CommandHandler, MessageHandler, filters)

from conversation_bot.login_manager import LoginHandler
from mongo import Subscription

# Define conversation states
CAMS_SELECTION, CAM_SELECT_ACTION, SETTINGS_INTERVAL, SETTINGS_DAILY, CONFIRM_SUBSCRIPTION = range(5)


class CamSubHandler:
    def __init__(self, get_by_user_id: callable):
        self.get_by_user_id = get_by_user_id

    async def handle_cams(self, update: Update, context: CallbackContext) -> int:
        # Retrieve user_id from context
        user_id = context.user_data.get('user_id')
        # Retrieve the cameras belonging to the user
        cams_id = self.get_by_user_id(user_id)
        if not cams_id:
            # End the conversation if the user has no cameras
            await update.message.reply_text('You don\'t have any cameras.')
            return ConversationHandler.END

        # Display the list of cameras for the user to choose from
        keyboard = [[f"/cam_choice {cam_id}" for cam_id in cams_id], ['/cancel']]
        reply_markup = ReplyKeyboardMarkup(keyboard, one_time_keyboard=True)
        await update.message.reply_text('Choose a camera:', reply_markup=reply_markup)

        return CAMS_SELECTION

    @staticmethod
    async def handle_cam_selection(update: Update, context: CallbackContext) -> int:
        # Save the chosen camera ID
        context.user_data['chosen_cam'] = context.args[0]
        chosen_cam_id = context.user_data.get('chosen_cam')
        # Check if the user is already subscribed to the chosen camera
        sub = Subscription.get_by_chat_and_cam_id(update.message.chat_id, chosen_cam_id)

        if sub and sub.is_subscribed:
            # User is already subscribed
            keyboard = [['/unsubscribe']]
            await update.message.reply_text(f'You are already subscribed to {chosen_cam_id}.')
        else:
            # User is not subscribed
            keyboard = [['/subscribe']]
        keyboard.append(['/cancel'])
        reply_markup = ReplyKeyboardMarkup(keyboard, one_time_keyboard=True)
        await update.message.reply_text('Choose an action:', reply_markup=reply_markup)
        return CAM_SELECT_ACTION

    @staticmethod
    async def handle_unsubscribe(update: Update, context: CallbackContext) -> int:
        # Determine the user's chosen action (subscribe or unsubscribe)
        command = update.message.text.split()[0].lower()
        chosen_cam_id = context.user_data.get('chosen_cam')
        # Create a Subscription instance for the chosen camera
        sub = Subscription(id=Subscription.Id(chat_id=update.message.chat_id, cam_id=chosen_cam_id))
        sub.is_subscribed = False
        sub.save()
        await update.message.reply_text(f'You are now {command}d from {sub.id.cam_id}.')
        # Log the subscription action
        logging.info("User %s %s from %s.", update.message.from_user.first_name, command, sub.id.cam_id)
        # Return to the main menu
        await LoginHandler.go_to_logged(update)

        return ConversationHandler.END

    @staticmethod
    async def handle_cancel(update: Update, _: CallbackContext) -> int:
        # Cancel the conversation and return to the main menu
        await LoginHandler.go_to_logged(update)
        return ConversationHandler.END

    @staticmethod
    async def handle_settings(update: Update, _: CallbackContext) -> int:
        # Ask the user for notification frequency
        await update.message.reply_text("How often do you want to be notified (in minutes)?",
                                        reply_markup=ReplyKeyboardRemove())
        return SETTINGS_INTERVAL

    @staticmethod
    async def handle_setting_interval(update: Update, context: CallbackContext) -> int:
        # Validate and process the user's input for notification interval
        notification_interval_str = update.message.text.strip()
        if not notification_interval_str.isdigit():
            # Invalid input, ask again
            await update.message.reply_text("Please enter a valid number for notification interval.")
            return SETTINGS_INTERVAL

        notification_interval = int(notification_interval_str)
        # Save the notification interval in user_data
        context.user_data['notification_interval'] = notification_interval

        # Ask the user for the daily time interval
        await update.message.reply_text("Specify the daily time interval (e.g., 22:00-06:00):")
        return SETTINGS_DAILY

    @staticmethod
    async def handle_setting_period(update: Update, context: CallbackContext) -> int:
        # Validate and process the user's input for daily time interval
        daily_time_interval = update.message.text.strip()
        if not re.match(r'^\d{2}:\d{2}-\d{2}:\d{2}$', daily_time_interval):
            # Invalid input format, ask again
            await update.message.reply_text("Please enter the daily time interval in the format HH:mm-HH:mm.")
            return SETTINGS_DAILY

        # Save the daily time interval in user_data
        context.user_data['daily_time_interval'] = daily_time_interval

        # Visualize user-provided settings and ask for confirmation
        notification_interval = context.user_data.get('notification_interval')
        daily_time_interval = context.user_data.get('daily_time_interval')
        chosen_cam_id = context.user_data.get('chosen_cam')
        confirmation_message = (f"Notification settings:\n"
                                f"Frequency: {notification_interval} minutes\n"
                                f"Daily time interval: {daily_time_interval}\n"
                                f"Camera: {chosen_cam_id}\n"
                                f"Type /subscribe to confirm or /cancel to go back.")
        reply_markup = ReplyKeyboardMarkup([["/subscribe", "/cancel "]], one_time_keyboard=True)
        await update.message.reply_text(confirmation_message, reply_markup=reply_markup)

        return CONFIRM_SUBSCRIPTION

    @staticmethod
    async def handle_confirm_subscription(update: Update, context: CallbackContext) -> int:
        # Determine the user's chosen action (subscribe or unsubscribe)
        chosen_cam_id = context.user_data.get('chosen_cam')
        notification_interval = context.user_data.get('notification_interval')
        daily_time_interval = context.user_data.get('daily_time_interval')
        # Create a Subscription instance for the chosen camera
        sub = Subscription(id=Subscription.Id(chat_id=update.message.chat_id, cam_id=chosen_cam_id))
        # Update subscription status and save to the database
        sub.is_subscribed = True
        sub.specs = Subscription.Specs(interval=notification_interval, period=daily_time_interval)
        sub.save()
        await update.message.reply_text(f'You are now subscribed to {sub.id.cam_id}.')
        # Log the subscription action
        logging.info("User %s subscribed to %s.", update.message.from_user.first_name, sub.id.cam_id)
        # Return to the main menu
        await LoginHandler.go_to_logged(update)
        return ConversationHandler.END

    def cam_sub_conv_handler(self) -> ConversationHandler:
        # Define the ConversationHandler for camera subscription
        return ConversationHandler(entry_points=[CommandHandler('cams', self.handle_cams)], states={
            CAMS_SELECTION: [CommandHandler('cam_choice', CamSubHandler.handle_cam_selection)],
            CAM_SELECT_ACTION: [CommandHandler('unsubscribe', CamSubHandler.handle_unsubscribe),
                                CommandHandler('subscribe', CamSubHandler.handle_settings)],
            SETTINGS_INTERVAL: [MessageHandler(filters.TEXT, CamSubHandler.handle_setting_interval)],
            SETTINGS_DAILY: [MessageHandler(filters.TEXT, CamSubHandler.handle_setting_period)],
            CONFIRM_SUBSCRIPTION: [CommandHandler('subscribe', CamSubHandler.handle_confirm_subscription)]},
                                   fallbacks=[CommandHandler('cancel', CamSubHandler.handle_cancel)],
                                   allow_reentry=True)
