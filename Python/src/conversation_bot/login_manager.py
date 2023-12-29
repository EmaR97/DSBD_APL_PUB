import logging

from telegram import Update, ReplyKeyboardRemove, ReplyKeyboardMarkup
from telegram.ext import CallbackContext, ConversationHandler, CommandHandler, MessageHandler, filters

from mongo import Chat

LOGGED, USERNAME, PASSWORD = range(3)


class LoginHandler:

    def __init__(self, cam_sub_conv_handler, login_request_handler):
        self.login_request_handler = login_request_handler
        self.cam_sub_conv_handler = cam_sub_conv_handler

    @staticmethod
    async def start_login(update: Update, context: CallbackContext) -> int:
        """
        Entry point for the login conversation. Checks if the user is already logged in
        and directs them to the appropriate state.
        """
        logging.info("login_start")
        chat = Chat.get_by_id(update.message.chat_id)
        if chat:
            context.user_data['user_id'] = chat.user_id
            await update.message.reply_text(f'Welcome back {chat.user_id}.\n')
            await LoginHandler.go_to_logged(update)
            return LOGGED
        await update.message.reply_text('Please enter your username.\n'
                                        'Send /cancel to stop talking to me.\n\n')
        return USERNAME

    @staticmethod
    async def handle_username(update: Update, context: CallbackContext) -> int:
        """
        Handles the user input for the username and checks if it exists in the database.
        """
        username = update.message.text.lower()
        logging.info(f"login_username:{username}")

        if username == '':
            await update.message.reply_text('Invalid input. Please click /start to start over.')
            return ConversationHandler.END
        context.user_data['username'] = username
        await update.message.reply_text('Please enter your password.')
        return PASSWORD

    async def handle_password(self, update: Update, context: CallbackContext) -> int:
        """
        Handles the user input for the password and completes the login process.
        """
        username = context.user_data.get('username')
        password = update.message.text.lower()

        # Do not display the password in the chat
        await update.message.delete()
        if not self.login_request_handler(username, password):
            await update.message.reply_text('Invalid input. Please click /start to start over.')
            return ConversationHandler.END
        # Save the chat ID upon completing the login
        chat_id = update.effective_message.chat_id
        chat = Chat(_id=chat_id, user_id=username)
        chat.save()
        await update.message.reply_text(f'Welcome {username}! You are now logged in.')
        await LoginHandler.go_to_logged(update)
        return LOGGED

    @staticmethod
    async def handle_logout(update: Update, context: CallbackContext) -> int:
        """
        Handles the user request to log out, cancels the conversation, and clears user data.
        """
        user = update.message.from_user
        logging.info("User %s canceled the conversation.", user.first_name)
        chat_id = update.effective_message.chat_id
        chat = Chat.get_by_id(chat_id)
        chat.delete()
        context.user_data.clear()
        await update.message.reply_text("Goodbye! I hope we can talk again some day.",
                                        reply_markup=ReplyKeyboardRemove())
        return ConversationHandler.END

    @staticmethod
    async def go_to_logged(update):
        """
        Displays the main menu options for a logged-in user.
        """
        keyboard = [['/cams', '/logout']]
        reply_markup = ReplyKeyboardMarkup(keyboard, one_time_keyboard=True)
        await update.message.reply_text('Choose an action:', reply_markup=reply_markup)

    def login_conv_handler(self) -> ConversationHandler:
        """
        Defines the ConversationHandler for the login process.
        """
        return ConversationHandler(entry_points=[CommandHandler('start', LoginHandler.start_login)], states={
            USERNAME: [MessageHandler(filters.TEXT & ~filters.COMMAND, LoginHandler.handle_username)],
            PASSWORD: [MessageHandler(filters.TEXT & ~filters.COMMAND, self.handle_password)],
            LOGGED: [self.cam_sub_conv_handler], }, fallbacks=[CommandHandler("logout", LoginHandler.handle_logout)],
                                   allow_reentry=True)
